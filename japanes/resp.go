package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Response 制定 HTTP Response 格式
type Response struct {
	r     *http.Request
	w     http.ResponseWriter
	Error *MyErr
	Data  map[string]interface{}
}

// NewResponse NewResponse
func NewResponse(w http.ResponseWriter, r *http.Request) *Response {
	resp := &Response{}
	resp.Data = make(map[string]interface{})
	resp.Error = &MyErr{
		ExtraInfo: make(map[string]interface{}),
	}
	resp.r = r
	resp.w = w

	return resp
}

// HTTPResponse 制定 HTTP Response 格式及方法
func (resp *Response) HTTPResponse() {
	JSONData, err := json.Marshal(resp)
	if err != nil {
		log.Printf("json.Marshal(resp) ERROR::: %+v, %+v\n", err, resp)
	}

	resp.w.Header().Set("Content-Type", "application/json")
	log.Printf("RESPONSE DATA::: %+v\n", string(JSONData))
	_, err = resp.w.Write(JSONData)
	if err != nil {
		log.Printf("RESPONSE WRITE ERROR::: %+v\n", err)
	}
}

// SetError SetError
func (resp *Response) SetError(err *MyErr) {
	resp.Error = err
}

func (resp *Response) addData(key string, value interface{}) {
	fmt.Println(value)
	resp.Data[key] = value
}
