package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"regexp"
	"runtime"
	"strconv"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func GetSliceIntFromBytes(input []byte) []TableId {
	sizeofInt := 8
	data := make([]TableId, len(input)/sizeofInt)
	buf := bytes.NewBuffer(input)
	for i := range data {
		var re int64
		binary.Read(buf, binary.LittleEndian, &re)
		data[i] = TableId(re)
	}

	return data
}

func GetSliceBytesFromInts(input []TableId) []byte {
	buf := new(bytes.Buffer)

	for i := range input {
		binary.Write(buf, binary.LittleEndian, int64(input[i]))
	}
	return buf.Bytes()
}

func ContainsIntSlice(s []TableId, e TableId) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (tId TableId) ToBytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, int64(tId))
	return buf.Bytes()
}

func (tId TableId) ToString() (str string) {
	str = strconv.FormatInt(int64(tId), 10)
	return
}

func (tId TableId) ToInt() int64 {
	return int64(tId)
}

func ToInt64FromBytes(st []byte) int64 {
	buf := bytes.NewReader(st)
	var result int64
	binary.Read(buf, binary.LittleEndian, &result)
	return result
}

//proc string like "createcustomernew(1)", and return 1
func GetIdFromStr(input string) TableId {
	valId := regexp.MustCompile(`(\d+)`)
	val, _ := strconv.Atoi(valId.FindString(input))
	return TableId(val)
}

func RetrieveByMapLevel(f interface{}, levels []string) interface{} {
	length := len(levels)

	if length <= 0 {
		return f
	}

	result := f.(map[string]interface{})

	for k, v := range result {
		if k == levels[0] {
			return RetrieveByMapLevel(v.(interface{}), levels[1:])
		}
	}

	return result
}

func HandleDecodeData(f map[string]interface{}) {
	for k, v := range f {
		log.Println("key", k)
		switch vv := v.(type) {
		//		case string:
		//			log.Println(k, "is string", vv)
		//		case int:
		//			log.Println(k, "is int", vv)
		//		case []interface{}:
		//			log.Println(k, "is an array:")
		//			for i, u := range vv {
		//				log.Println(i, u)
		//			}
		//		case interface{}:
		//			log.Println(k, "is an interface")
		//			HandleDecodeData(f[k].(map[string]interface{}))
		default:
			log.Println(vv)
			log.Println(k, "is of a type I don't know how to handle")
		}
	}
}

func HandleError(err error) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	log.Println(err, file, "line:", line)
}
