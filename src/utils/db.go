package utils

import (
	"fmt"
	"io"

	driver "github.com/johnnadratowski/golang-neo4j-bolt-driver"
	"github.com/johnnadratowski/golang-neo4j-bolt-driver/structures/graph"
)

// URI sent as Environment variable (secret) or defaults localinstance call
const (
	URI = "bolt://neo4j:test@neo4j:7687"
)

func CreateConnection() driver.Conn {
	driver := driver.NewDriver()
	con, err := driver.OpenNeo(URI)
	handleError(err)
	return con
}

// Here we create a simple function that will take care of errors, helping with some code clean up
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

// Here we prepare a new statement. This gives us the flexibility to
// cancel that statement without any request sent to Neo
func PrepareSatement(query string, con driver.Conn) driver.Stmt {
	st, err := con.PrepareNeo(query)
	handleError(err)
	return st
}
func QueryStatement(st driver.Stmt) driver.Rows {
	// Even once I get the rows, if I do not consume them and close the
	// rows, Neo will discard and not send the data
	rows, err := st.QueryNeo(nil)
	handleError(err)
	return rows
}

func consumeMetadata(rows driver.Rows, st driver.Stmt) {
	// Here we loop through the rows until we get the metadata object
	// back, meaning the row stream has been fully consumed
	defer st.Close()
	var err error
	err = nil

	for err == nil {
		var row []interface{}
		row, _, err = rows.NextNeo()
		if err != nil && err != io.EOF {
			panic(err)
		} else if err != io.EOF {
			fmt.Printf("PATH: %#v\n", row[0].(graph.Path)) // Prints all paths
		}
	}
}

func ConsumeRows(rows driver.Rows, st driver.Stmt) {
	// This interface allows you to consume rows one-by-one, as they
	// come off the bolt stream. This is more efficient especially
	// if you're only looking for a particular row/set of rows, as
	// you don't need to load up the entire dataset into memory
	for {
		data, _, err := rows.NextNeo()
		if err != nil && err == io.EOF {
			break
		}
		handleError(err)
		fmt.Printf("COLUMNS: %#v\n", rows.Metadata()["fields"].([]interface{})) // COLUMNS: n.foo,n.bar
		fmt.Printf("FIELDS: %d %f\n", data[0].(int64), data[1].(float64))
	}

	// This query only returns 1 row, so once it's done, it will return
	// the metadata associated with the query completion, along with
	// io.EOF as the error
	//_, _, err := rows.NextNeo()
	// FIELDS: 1 2.2

	st.Close()
}

// Executing a statement just returns summary information
func ExecuteStatement(st driver.Stmt, props map[string]interface{}) driver.Result {
	tmp := map[string]interface{}{"props": props}
	result, err := st.ExecNeo(tmp)
	handleError(err)
	numResult, err := result.RowsAffected()
	handleError(err)
	fmt.Printf("CREATED ROWS: %d\n", numResult) // CREATED ROWS: 1
	return result
}
