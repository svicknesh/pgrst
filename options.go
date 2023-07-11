package pgrst

import (
	"strconv"
	"strings"
)

type FilterType uint8
type CountType uint8

const (
	Equal FilterType = iota + 1
	GreaterThan
	GreaterThanEqual
	LessThan
	LessThanEqual
	NotEqual
)

const (
	CountTypeExact FilterType = iota + 1
	CountTypePlanned
	CountTypeEstimated
)

// Select - specify fields to be returned from the table
func (db *DB) Select(col ...string) *DB {
	db.SetQuery("select", strings.Join(col, ","))
	return db
}

// Limit - limits the returned results
func (db *DB) Limit(limit int) *DB {
	db.SetQuery("limit", strconv.Itoa(limit))
	return db
}

// Offset - returns result from the given offset
func (db *DB) Offset(offset int) *DB {
	db.SetQuery("offset", strconv.Itoa(offset))
	return db
}

// WithCount - sets the correct header to get actual count of items from the database
func (db *DB) WithCount(ct CountType) *DB {
	db.SetHeader("prefer", "count="+ct.String())
	return db
}

// SetHeader - sets custom header for PostgREST to support features not immediately implemented by this library
func (db *DB) SetHeader(key, value string) *DB {
	db.client.SetHeader(key, value)
	return db
}

// SetQuery - sets custom query parameters for PostgREST to support features not immediately implemented by this library
func (db *DB) SetQuery(key, value string) *DB {
	db.qVal.Add(key, value)
	return db
}

// Where - filter result rows by adding conditions on columns
func (db *DB) Where(field string, ft FilterType, value string) *DB {
	db.qVal.Add(field, ft.String()+"."+value)
	return db
}

// WhereIn - filter result rows by adding conditions on columns, one of a list of values
func (db *DB) WhereIn(field string, values []string) *DB {
	db.qVal.Add(field, "in.("+strings.Join(values, ",")+")")
	return db
}

func (ft FilterType) String() (str string) {
	ftStr := []string{"", "eq", "gt", "gte", "lt", "lte", "neq"} // entry 0 is always empty
	return ftStr[int(ft)]
}

func (ct CountType) String() (str string) {
	ctStr := []string{"", "exact", "planned", "estimated"} // entry 0 is always empty
	return ctStr[int(ct)]
}
