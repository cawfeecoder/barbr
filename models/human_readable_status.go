package models

type HumanReadableStatus struct {
	Type string `json:"type"`
	Message string `json:"message"`
	Param string `json:"param"`
	Value interface{} `json:"value"`
	Source string `json:"source"`
}