# Golang PostgREST library

Golang library to communicate with an existing PostgREST instance. This library crafts the necessary 
request structure needed.

This library implements the most commonly used features. Oher features will be added as time goes by.


## Example

### Creating new instance of database
```go

// the address can be direct access to the `PostgREST` or behind a reverse proxy.
address := "http://example.com:3000"
role := "webuser"
secretkey := "supersecretkey"
duration := time.Duration(time.Second*30)

db, err := pgrst.New(address, role, secretkey, duration)
if nil!=err{
    log.Println(err)
    os.Exit(1)
}

```

the parameters for creating a new instance of db are as follows

| Variables | Description | 
| :-- | :-- |
| `address` | Address of the postgrest instance, this could be directly exposed on a local network or behind a reverse proxy on the public internet. |
| `role` | Name of the role whihc has permission to act on the database. |
| `secretkey` | The key used to create a signed request sent to the postgrest instance to prevent unauthorized accesss. |
| `duration` | Maximum validity of the signed request, recommend keeping it between 30 - 45 seconds. |
|||


### Select 
```go
// users - slice of []User structure to store multiple values
// the returned values are mapped to the structure using the json field name
err := db.NewSelect().Table("users").Limit(25).Offset(0).Output(&users).Exec()
```

```go
// users - slice of []User structure to store multiple values
// the returned values are mapped to the structure using the json field name
// the where condition requires the key name to compare and the value, even for an numeric value, it must be sent as a string
dbSelect := db.NewSelect().Table("users").Where("id",pgrst.Equal,"1").Output(&users)
err := dbSelect.Exec()
if nil != err {
    log.Println(err)
    os.Exit(1)
}

if !dbSelect.IsSuccess() {
    log.Println("no result found with given criteria")
    os.Exit(1)
}

```


### Insert 
```go
// user - structure storing user details
// the json field of user must match the table field names
err := db.NewInsert().Table("users").Input(user).Exec()
```


### Update 
```go
// similar to insert but requires a where condition to update
// the where condition requires the key name to compare and the value, even for an numeric value, it must be sent as a string
err := db.NewUpdate().Table("users").Where("id", pgrst.Equal,"1").Input(user).Exec()
```


### Call stored procedure to select data 
```go
// usersaddrs - slice of []UserAddr structure to store multiple values
// this calls the 
// the returned values are mapped to the structure using the json field name
err := db.NewSelect().RPC("getuserwithaddress").SetQuery("email", "user@example.com").Output(&usersaddrs).Exec()
```


### Call stored procedure to insert data 
```go
err := db.NewSelect().RPC("setuseraddress").Input(usersaddrs).Exec()
```


## Constants

The following constants are used for `Where` conditions
| Constants | Description | 
| :-- | :-- |
| `Equal` | Search for values that match exactly the input. |
| `GreaterThan` | Search for values greater than the input. Used for numberic values. |
| `GreaterThanEqual` | Search for values greater or equal the input. Used for numberic values. |
| `LessThan` | Search for values less than the input. Used for numberic values. |
| `LessThanEqual` | Search for values less than the or equal the input. Used for numberic values. |
| `NotEqual` | Search for values that don't match the input. |
|||

The following constants are used for `WithCount` conditions
| Constants | Description | 
| :-- | :-- |
| `CountTypeExact` | Returns the total size of the table based on the query and parameters. This could be slower to run on larger databases. |
| `CountTypePlanned` | Returns the total size of the table based on the query and parameters using PostgreSQL statistics. Fairly accurate and fast. |
| `CountTypeEstimated` | Similar as above but uses exact count until a threshold and get the planned count after. |
|||

## Functions

The following functions are available in this library

| Function | Description |
| :-- | :-- |
| `New` | creates a new instance of database with the `address`, `role`, `secretkey` and `duration`. |
| `NewSelect` | Sets the request method to `GET` to search through the database. Used with `Table` and `RPC`. |
| `NewInsert` | Sets the request method to `POST` to insert data to the database.  Used with `Table` and `RPC`. |
| `NewUpsert` | Sets the request method to `POST` to update if exist otherwise update data in the database. Sets the necessary headers requird by `PostgREST`. |
| `NewUpdate` | Sets the request method to `PATCH` to update an existing data in the database. |
| `NewDelete` | Sets the request method to `DELETE` to delete items from the database. |
| `SetRole` | overrides the `role` configured in `New`.  |
| `Input` | sets the input for inserting, upserting, or updating data. The `json` fields **MUST** match the column names in the related tables.  |
| `Output` | sets the pointer to a structure to store the result. The type can be a `slice`, `struct` or `string`. |
| `Debug` | prints the request being sent and the response received to the console. |
| `Exec` | Executes the chained requet.  |
| `Select` | define the fields to be returned from the `NewSelect` request.  |
| `Where` | sets the condition for the request. Refer to teh `Constants` section for futher explanation.  |
| `WhereIn` | like above but compares against a list of values. |
| `Limit` | number of results to be returned. |
| `Offset` | skips that many rows, used with `Limit`. |
| `WithCount` | returns number of items in a given `Table`. Refer to the `Constants` section for further explanation.  |
| `SetQuery` | sets custom URL query parameters, allowing support of other `PostgREST` features. |
| `SetHeader` | sets custom HTTP headers, allowing support of other `PostgREST` features. |
| `IsSuccess` | determines if the request is successful. Requests with HTTP code *2xx* are considered successful.  |
| `GetCount` | Gets the total count of items in the given `Table`. Requires `WithCount`. |
| `GetRange` | Gets the range of the items for this request. |
|||

