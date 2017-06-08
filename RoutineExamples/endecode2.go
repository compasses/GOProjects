package main

import (
	"time"
	"encoding/json"
	"fmt"
	"errors"
)

type Dog struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Breed string `json:"breed"`
	BornAt Time `json:"born_at"`
}

type Time struct {
	time.Time
}

func (t Time) MarshalJSON()([]byte, error)  {
	return json.Marshal(t.Time.Unix())
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var i int64
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	t.Time = time.Unix(i, 0)
	return nil
}

//type JSONDog struct {
//	Dog
//	BornAt int64 `json:"born_at"`
//}

//func NewJSONDog(dog Dog) JSONDog {
//	return JSONDog{
//		dog,
//		dog.BornAt.Unix(),
//	}
//}
//
//func (jd JSONDog) ToDog() Dog {
//	return Dog {
//		jd.Dog.ID,
//		jd.Dog.Name,
//		jd.Dog.Breed,
//		time.Unix(jd.BornAt, 0),
//	}
//}

func test2() {
	dog := Dog{1, "bowser", "husky", Time{time.Now()}}
	b, err := json.Marshal(dog)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))

	b = []byte(`{
    "id":1,
    "name":"bowser",
    "breed":"husky",
    "born_at":1480979203}`)
	dog = Dog{}
	json.Unmarshal(b, &dog)
	fmt.Println(dog)
}

type BankAccount struct {
	ID            string `json:"id"`
	Object        string `json:"object"`
	RoutingNumber string `json:"routing_number"`
}

type Card struct {
	ID     string `json:"id"`
	Object string `json:"object"`
	Last4  string `json:"last4"`
}

type Data struct {
	*Card
	*BankAccount
}

func (d Data) MarshalJSON() ([]byte, error) {
	if d.Card != nil {
		return json.Marshal(d.Card)
	} else if d.BankAccount != nil {
		return json.Marshal(d.BankAccount)
	} else {
		return json.Marshal(nil)
	}
}

func (d* Data) UnmarshalJSON(data []byte) error {
	temp := struct {
		Object string `json:"object"`
	}{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	if temp.Object == "card" {
		var c Card
		if err := json.Unmarshal(data, &c); err != nil {
			return err
		}
		d.Card = &c
		d.BankAccount = nil
		return nil
	} else if temp.Object == "bank_account" {
		var b BankAccount
		if err := json.Unmarshal(data, &b); err != nil {
			return err
		}
		d.BankAccount = &b
		d.Card = nil
		return nil
	} else {
		return errors.New("Invalid object value")
	}

	return nil
}

func test3() {
	jsonStr := `
{
  "data": {
    "object": "card",
    "id": "card_123",
    "last4": "4242"
  }
}
`
	var m map[string]map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &m); err != nil {
		panic(err)
	}
	fmt.Println(m)

	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
func main() {
	//test2()
	//test3()
	jsonStr := `
		{
		"data": {
			"object": "card",
			"id": "card_123",
			"last4": "4242"
			}
		}
		`
	var m map[string]Data
	if err := json.Unmarshal([]byte(jsonStr), &m); err != nil {
		panic(err)
	}

	fmt.Println(m)
	data := m["data"]
	if data.Card != nil {
		fmt.Println(data.Card)
	}
	if data.BankAccount != nil {
		fmt.Println(data.BankAccount)
	}

	b, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
}
