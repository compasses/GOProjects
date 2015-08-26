package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"math/rand"
	//"os"
	//"time"
	//"encoding/binary"
	"strconv"
	//"encoding/json"
	"github.com/boltdb/bolt"
	"log"
)

var GlobalDB *bolt.DB

func RepoCreateATSRsp(req *ATSReq) []ATSRsp {
	var rest []ATSRsp

	for _, atsR := range req.SkuIds {
		rsp := ATSRsp{
			SkuId:          atsR,
			Ats:            10,
			AllowBackOrder: true,
		}
		rest = append(rest, rsp)
	}

	log.Printf("ATS Rsp %+v\n", rest)

	return rest
}

func RepoCreateRecommandationProducts(key []byte, Id int) []int {
	var ProductId []int

	GlobalDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(key)
		if err != nil {
			return err
		}

		all := b.Get([]byte("Products"))
		if all != nil {
			ProductId = getSliceIntFromBytes(all)
		}

		if !containsIntSlice(ProductId, Id) {
			ProductId = append(ProductId, Id)
		}

		b.Put([]byte("Products"), getSliceBytesFromInts(ProductId))

		return nil
	})

	return ProductId
}
func RepoCreateAccount(key []byte, customer CustomerCreate) CustomerCreateRsp {
	var result CustomerCreateRsp
	key = append(key, []byte("ACCOUNT")...)

	GlobalDB.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(key)
		if err != nil {
			return err
		}
		user := b.Get([]byte(customer.Account))

		if user == nil {
			//create new user
			customer.AccountInfo.CustomerID = (tx.Size() * rand.Int63()) % 997
			customer.AccountInfo.CustomerCode = "offline" + strconv.FormatInt(customer.AccountInfo.CustomerID, 10)
			customer.AccountInfo.ChannelAccountID = customer.ChannelId * customer.AccountInfo.CustomerID
			cusStream, _ := json.Marshal(&customer)
			b.Put([]byte(customer.Account), cusStream)
			//use Id to get account
			buf := new(bytes.Buffer)
			binary.Write(buf, binary.LittleEndian, customer.AccountInfo.CustomerID)
			b.Put(buf.Bytes(), []byte(customer.Account))
			result = customer.AccountInfo
		}

		json.Unmarshal(user, &customer)
		result = customer.AccountInfo
		return nil
	})
	return result
}

func encode(v int) []byte {
	var result []byte
	result = strconv.AppendInt(result, int64(v), 10)
	return result
}

func getSliceIntFromBytes(input []byte) []int {
	sizeofInt := 4
	data := make([]int, len(input)/sizeofInt) // int used 4 bytes
	for i := range data {
		num, _ := binary.Varint(input[i*sizeofInt : (i+1)*sizeofInt])
		data[i] = int(num)
	}
	return data
}

func getSliceBytesFromInts(input []int) []byte {
	result := make([]byte, len(input)*4)
	for i := range input {
		binary.PutVarint(result[i*4:], int64(input[i]))
	}
	return result
}

func containsIntSlice(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func init() {
	GlobalDB, _ = bolt.Open("./EshopOfflineServerDB", 0666, nil)
//	go func() {
//		prev := GlobalDB.Stats()
//		for {
//			time.Sleep(10 * time.Second)

//			states := GlobalDB.Stats()
//			diff := states.Sub(&prev)
//			//json.NewEncoder(os.Stderr).Encode(diff)
//			prev = states
//		}
//	}()
}
