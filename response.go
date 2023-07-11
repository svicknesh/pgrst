package pgrst

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
	return &db.dbResponse
}
