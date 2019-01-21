package models

import "encoding/json"

type Response struct {
	Data []interface{} `json:"data"`
	Errors []interface{} `json:"errors"`
}

func (res *Response) ToJSON() (data []byte, err error){
	data, err = json.Marshal(res)
	return
}