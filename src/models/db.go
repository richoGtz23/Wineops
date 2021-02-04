// TODO: official NEO4j https://medium.com/@angadsharma1016/optimizing-go-neo4j-concurrency-patterns-810dff25f88f
package models

import (
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type NeoDb struct {
	neo4j.Driver
}

func CreateConnection(neo4jHost, neo4jPassword, neo4jUser, neo4jPort string) *NeoDb {
	configForNeo4j40 := func(conf *neo4j.Config) {
		conf.Log = neo4j.ConsoleLogger(neo4j.DEBUG)
	}
	uri := fmt.Sprintf("bolt://%s:%s", neo4jHost, neo4jPort)
	fmt.Println(uri)
	d, err := neo4j.NewDriver(uri, neo4j.BasicAuth(neo4jUser, neo4jPassword, ""), configForNeo4j40)
	handleError(err)
	return &NeoDb{d}
}
func CreateSession(d *NeoDb, mode string) neo4j.Session {
	m := neo4j.AccessModeRead
	if mode == "write" {
		m = neo4j.AccessModeWrite
	}
	sessionConfig := neo4j.SessionConfig{AccessMode: m}
	session := d.NewSession(sessionConfig)
	return session
}

// Here we create a simple function that will take care of errors, helping with some code clean up
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func DefaultTxMetadata(user string) neo4j.TransactionConfig {
	return neo4j.TransactionConfig{
		Metadata: map[string]interface{}{"user": user, "datetime": time.Now()},
	}
}
