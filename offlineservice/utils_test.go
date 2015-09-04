package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"regexp"
	"strconv"
	"testing"
)

func TestBytesToInt(t *testing.T) {
	test1 := int64(256)
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, test1)
	//binary.PutVarint(buf.Bytes(), int64(test1))
	fmt.Println("byte is ", buf.Bytes())
	bufr := bytes.NewBuffer(buf.Bytes())
	var rest int64
	err := binary.Read(bufr, binary.LittleEndian, &rest)
	if err != nil {
		t.Error(err)
	}
	fmt.Println("rest value ", rest)

	num, _ := binary.Varint(buf.Bytes())
	fmt.Println("num is ", num)

	var testS []TableId
	testS = append(testS, 34, 3443, 344242)
	tids := GetSliceBytesFromInts(testS)
	fmt.Println("length is ", len(tids), "val is ", tids)

	tidr := GetSliceIntFromBytes(tids)
	fmt.Println("int value is ", tidr)

	t4 := ToInt64FromBytes(tids[0:8])
	fmt.Println("t4 val is ", t4)

}

func TestBytesToInt2(t *testing.T) {
	var intVar = 123
	fmt.Println("intVar is : ", intVar)

	intByte := []byte(strconv.Itoa(intVar))

	fmt.Println("intByte is : ", intByte)
}

func TestGetIntFromStr(t *testing.T) {
	input := "CreateCustomerNew(23)"
	valId := regexp.MustCompile(`(\d+)`)
	val, _ := strconv.Atoi(valId.FindString(input))
	fmt.Println("id :", val)
	var cc CustomerAddress
	var ni []interface{}
	ni = append(ni, "ok")
	ni = append(ni, 2323)
	ni = append(ni, cc)

	fmt.Println("interfaces values ", ni)
}
