package main

import (
	"github.com/compasses/GOProjects/RoutineExamples/protobuf/example"
	"github.com/golang/protobuf/proto"
	"log"
)

func main() {
	test1 := &test.Test{
		Label:         proto.String("hello"),
		Type:          proto.Int32(17),
		Reps:          []int64{1, 2, 3},
		Optionalgroup: &test.Test_OptionalGroup{RequiredField: proto.String("Good bye")},
	}

	data, err := proto.Marshal(test1)
    log.Printf("original data:%q\n", data)

	if err != nil {
		log.Fatal("marshaling error :", err)
	}
	newTest := &test.Test{}
	err = proto.Unmarshal(data, newTest)

	if err != nil {
		log.Fatal("Unmarshaling error:", err)
	}

	if test1.GetLabel() != newTest.GetLabel() {
		log.Fatal("data mismatch %q != %q", test1.GetLabel(), newTest.GetLabel())
	}

    log.Printf("get new data=%q\n", newTest)

	p := &test.Person{
		Id:    proto.Int32(1234),
		Name:  proto.String("Jet Doe"),
		Email: proto.String("je"),
		Phones: []*test.Person_PhoneNumber{
			{Number: proto.String("555-4321")},
		},
	}

	data, err = proto.Marshal(p)

	if err != nil {
		log.Fatal("marshaling error :", err)
	}

	log.Printf("data is %q\n", data)
}
