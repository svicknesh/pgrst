package pgrst

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/svicknesh/httpclient"
)

type httpMethod uint8

const (
	httpMethodGet httpMethod = iota + 1
	httpMethodPost
	httpMethodPatch
	httpMethodDelete
)

// DB - remote DB information
type DB struct {
	client            *httpclient.Request
	duration          time.Duration
	jsvo              jws.SignVerifyOption
	authz             authz
	method            httpMethod
	endpoint          string
	qVal              url.Values
	input, output     any
	success, debug    bool
	dbResponse        DBResponse
	start, end, total int
}

// authz - authorization token with remote server
type authz struct {
	Expiry int64  `json:"exp"`
	Role   string `json:"role"`
}

// New - creates new instance of Database
func New(address, role, secret string, duration time.Duration) (db *DB, err error) {

	db = new(DB)

	db.client = httpclient.NewRequest(address, time.Duration(time.Second*30), &tls.Config{InsecureSkipVerify: true}, nil)
	db.client.SetUserAgent("pgrst-golang/v1.0")

	secretKey, err := jwk.FromRaw([]byte(secret))
	if nil != err {
		return nil, fmt.Errorf("newdb: jwk raw error -> %w", err)
	}

	db.qVal = make(url.Values)
	db.duration = duration

	db.jsvo = jws.WithKey(jwa.HS256, secretKey)

	db.authz.Role = role

	return
}

// genToken - creates an auth token with the given role for request
func (db *DB) genToken() (token string) {

	db.authz.Expiry = time.Now().Add(db.duration).UTC().Unix()

	bytes, _ := json.Marshal(db.authz)
	buf, _ := jws.Sign(bytes, db.jsvo)

	return string(buf)
}

// NewSelect - returns a new instance for selecting data
func (db *DB) NewSelect() *DB {
	clone := *db
	clone.method = httpMethodGet
	clone.qVal = make(url.Values) // reset this to a new instance
	return &clone
}

// NewInsert - returns a new instance for inserting data
func (db *DB) NewInsert() *DB {
	clone := *db
	clone.method = httpMethodPost
	clone.qVal = make(url.Values) // reset this to a new instance
	return &clone
}

// NewUpsert - returns a new instance for inserting or updating data (if exists)
func (db *DB) NewUpsert() *DB {
	clone := *db
	clone.method = httpMethodPost
	clone.qVal = make(url.Values) // reset this to a new instance
	clone.client.SetHeader("resolution", "merge-duplicates")
	return &clone
}

// NewUpdate - returns a new instance for updating data
func (db *DB) NewUpdate() *DB {
	clone := *db
	clone.method = httpMethodPatch
	clone.qVal = make(url.Values) // reset this to a new instance
	return &clone
}

// NewDelete - returns a new instance for deleting data
func (db *DB) NewDelete() *DB {
	clone := *db
	clone.method = httpMethodDelete
	clone.qVal = make(url.Values) // reset this to a new instance
	return &clone
}

// RPC - executes a remote function on the database
func (db *DB) RPC(funcName string) *DB {
	db.endpoint = "/rpc/" + funcName + "?"
	return db
}

// Table - sets the name of the table to operate on
func (db *DB) Table(tableName string) *DB {
	db.endpoint = "/" + tableName + "?"
	return db
}

// SetRole - sets the name of the role to use for this request
func (db *DB) SetRole(role string) *DB {
	db.authz.Role = role
	return db
}

// Input - structure to send to PostgREST instance
func (db *DB) Input(input any) *DB {
	db.input = input
	return db
}

// Output - structure to save a single result from PostgREST instance, remember to pass a pointer to the 'output'
func (db *DB) Output(output any) *DB {
	db.output = output
	return db
}

// Debug - prints the HTTP Status code and any messages received to console
func (db *DB) Debug() *DB {
	db.debug = true
	return db
}

// Exec - executes this query
func (db *DB) Exec() (err error) {

	var reader io.Reader
	if nil != db.input {
		reader, err = NewReader(db.input)
		if nil != err {
			return fmt.Errorf("exec: create reader error -> %w", err)
		}
	}

	//fmt.Println(db.endpoint + db.qVal.Encode())
	db.client.SetHeader("Authorization", "Bearer "+db.genToken())

	var response *httpclient.Response

	switch db.method {
	case httpMethodGet:
		if nil == db.output {
			return fmt.Errorf("exec: no output destination specificed for SELECT")
		}
		response, err = db.client.Get(db.endpoint + db.qVal.Encode())

	case httpMethodPost:
		response, err = db.client.Post(db.endpoint+db.qVal.Encode(), reader)

	case httpMethodPatch:
		response, err = db.client.Patch(db.endpoint+db.qVal.Encode(), reader)

	case httpMethodDelete:
		response, err = db.client.Delete(db.endpoint + db.qVal.Encode())

	default:
		err = fmt.Errorf("exec: SELECT, INSERT, UPSERT, UPDATE or DELETE not specified")
	}

	if db.debug {
		log.Printf("debug -> HTTP Status Code (%d), HTTP Response -> %s", response.StatusCode, response.Buffer.String())
	}

	//fmt.Println(response)
	//fmt.Println("1", response.StatusCode, db.method, response.Buffer.String())

	if nil != err {
		return fmt.Errorf("exec: request error -> %w", err)
	}

	var contentRange string
	crs := response.GetHeader("content-range")
	if len(crs) == 1 {
		contentRange = crs[0] // we are only interested in the content range if there is 1 items in it
	}

	switch response.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusNoContent, http.StatusPartialContent:
		// request completed successfully, return
		// 200 and 206 is for SELECT, 201 is for INSERT, 204 is for UPDATE, UPSERT or DELETE

		switch db.method {

		case httpMethodPatch:
			if contentRange != "*/*" {
				db.success = true
			} else {
				// if the content range returns '*/*' it means the update was not successful, likely the update condition is incorrect
				db.success = false
			}
		default:
			//fmt.Println("2", response.StatusCode, db.method, response.Buffer.String())
			// other methods can be assumed to be successful
			db.success = true
		}

		// save the offset and count
		if len(contentRange) != 0 {
			fmt.Sscanf(contentRange, "%d-%d/%d", &db.start, &db.end, &db.total)
		}

		//fmt.Println(response.Buffer.String())
		err = setOutput(response.Buffer.Bytes(), db.output)
		if nil != err {
			return fmt.Errorf("error setting output: %w", err)
		}

	case http.StatusConflict:
		// this happens for INSERT when there is a conflicting rule, save the database response
		response.ToJSON(&db.dbResponse)
		db.dbResponse.HTTPStatusCode = response.StatusCode

		return

	default:
		// any other HTTP status return it as an error
		response.ToJSON(&db.dbResponse)
		db.dbResponse.HTTPStatusCode = response.StatusCode

		return fmt.Errorf("exec: error completing request -> %s", db.dbResponse.String())
	}

	return
}
