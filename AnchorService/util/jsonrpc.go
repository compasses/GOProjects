package util

import (
	"encoding/json"
	"fmt"
)

func EncodeJSON(data interface{}) ([]byte, error) {
	encoded, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return encoded, nil
}

func EncodeJSONString(data interface{}) (string, error) {
	encoded, err := EncodeJSON(data)
	if err != nil {
		return "", err
	}
	return string(encoded), err
}

type JSONError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewJSONError(code int, message string, data interface{}) *JSONError {
	j := new(JSONError)
	j.Code = code
	j.Message = message
	j.Data = data
	return j
}

func (e *JSONError) Error() string {
	s := fmt.Sprint(e.Message)
	if e.Data != nil {
		s += fmt.Sprint(": ", e.Data)
	}
	return s
}

type JSON2Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Params  json.RawMessage `json:"params,omitempty"`
	Method  string          `json:"method,omitempty"`
}

func NewJSON2Request(method string, id, params interface{}) *JSON2Request {
	j := new(JSON2Request)
	j.JSONRPC = "2.0"
	j.ID = id
	if b, err := json.Marshal(params); err == nil {
		j.Params = b
	}
	j.Method = method
	return j
}

func ParseJSON2Request(request string) (*JSON2Request, error) {
	j := new(JSON2Request)
	err := json.Unmarshal([]byte(request), j)
	if err != nil {
		return nil, err
	}
	if j.JSONRPC != "2.0" {
		return nil, fmt.Errorf("Invalid JSON RPC version - `%v`, should be `2.0`", j.JSONRPC)
	}
	return j, nil
}

func (j *JSON2Request) JSONString() (string, error) {
	return EncodeJSONString(j)
}

func (j *JSON2Request) String() string {
	str, _ := j.JSONString()
	return str
}

type JSON2Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Error   *JSONError      `json:"error,omitempty"`
	Result  json.RawMessage `json:"result,omitempty"`
}

func NewJSON2Response() *JSON2Response {
	j := new(JSON2Response)
	j.JSONRPC = "2.0"
	return j
}

func (j *JSON2Response) JSONString() (string, error) {
	return EncodeJSONString(j)
}

func (j *JSON2Response) JSONResult() []byte {
	return j.Result
}

func (j *JSON2Response) String() string {
	str, _ := j.JSONString()
	return str
}
