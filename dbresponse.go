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

// IsSuccess - indicates if the `Exec` request was a success
func (db *DB) IsSuccess() (success bool) {
	return db.success
}

// GetCount - returns the number of items found, for SELECT requests
func (db *DB) GetCount() (count int) {
	return db.total
}

// GetRange - returns the range of items found, for SELECT requests
func (db *DB) GetRange() (start, end int) {
	return db.start, db.end
}

// GetDBResponse - returns the database response
func (db *DB) GetDBResponse() (dbResponse *DBResponse) {
	return db.dbResponse
}
