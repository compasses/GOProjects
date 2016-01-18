package online

import (
	"log"
	"net/http"
	"reflect"
)

var SuccNum int  = 0
var FailNum int  = 0

func LogOutPut(NeedLog bool,v ...interface{}) {
	if NeedLog {
		log.Println(v)
	}
}

func RequstFormat(NeedLog bool, req *http.Request, newbody string) {
	if !NeedLog {
		return
	}
	result := "URL: " + req.URL.String() + "\r\n"
	result += "Method: " + req.Method + "\r\n"
	result += "Body: " + newbody + "\r\n"
	result += "Header: "
	for key, _ := range req.Header {
		var vals string = ""
		for _, allV := range req.Header[key] {
			vals += allV
		}
		result += " Key: " + key + " -> " + vals + "\r\n"
	}
	LogOutPut(NeedLog, result)
}

func ResponseFormat(NeedLog bool, resp *http.Response, body string) {
	if !NeedLog {
		return
	}
	
	result := "Status: " + resp.Status + "\r\n"
	result += "Body: " + body + "\r\n"
	result += "Header: "
	for key, _ := range resp.Header {
		var vals string = ""
		for _, allV := range resp.Header[key] {
			vals += allV
		}
		result += " Key: " + key + " -> " + vals + "\r\n"
	}
	LogOutPut(NeedLog, result)
}

func ReflectStruct(req *http.Request) {
	s := reflect.ValueOf(req).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		log.Printf("%d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}
}
