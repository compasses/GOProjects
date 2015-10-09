package online

import (
	"log"
	"os"
	"strings"
	"testing"

	"golang.org/x/text/encoding"
	//"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func load(direction string, enc encoding.Encoding) (func() transform.Transformer, error) {

	newTransformer := enc.NewEncoder
	if direction == "Decode" {
		newTransformer = enc.NewDecoder
	}

	return newTransformer, nil
}

func TestConversion(t *testing.T) {
	f, err := os.OpenFile("testlogfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		t.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Println("This is a test log entry")

	sr := strings.NewReader("你好，世界")

	newTransformer, _ := load("Decode", simplifiedchinese.GBK)

	rInUTF8 := transform.NewReader(sr, newTransformer())
	res := make([]byte, 100)
	rInUTF8.Read(res)

	log.Println("t is ", "真的吗", string(res))
}
