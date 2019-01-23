package models

import "errors"

type HumanReadableStatus struct {
	Type string `json:"type"`
	Message string `json:"message"`
	Param string `json:"param"`
	Value interface{} `json:"value"`
	Source string `json:"source"`
}

func ToErrorsArray(hr []HumanReadableStatus) (err []error){
	for _, val := range hr {
		err = append(err, errors.New(val.Message))
	}
	return err
}