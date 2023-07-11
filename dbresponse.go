package pgrst

import (
	"encoding/json"
)

// DBResponse - database response
type DBResponse struct {
	Code           string `json:"code,omitempty"`
	Details        string `json:"details,omitempty"`
	Hint           string `json:"hint,omitempty"`
	Message        string `json:"message,omitempty"`
	HTTPStatusCode int    `json:"httpStatusCode,omitempty"`
}

func (dbResponse *DBResponse) HasCode() (b bool) {
	return len(dbResponse.Code) != 0
}

func (dbResponse *DBResponse) String() (str string) {
	bytes, _ := json.Marshal(dbResponse)
	return string(bytes)
}
