// TODO: official NEO4j https://medium.com/@angadsharma1016/optimizing-go-neo4j-concurrency-patterns-810dff25f88f
package models

import (
	"time"

	"github.com/neo4j/neo4j-go-driver/neo4j"
)

// URI sent as Environment variable (secret) or defaults localinstance call
const (
	URI      = "bolt://neo4j:7687"
	USER     = "neo4j"
	PASSWORD = "test"
)

type NeoDb struct {
	neo4j.Driver
}

func CreateConnection() *NeoDb {
	configForNeo4j40 := func(conf *neo4j.Config) {
		conf.Log = neo4j.ConsoleLogger(neo4j.DEBUG)
		conf.Encrypted = false
	}
	d, err := neo4j.NewDriver(URI, neo4j.BasicAuth(USER, PASSWORD, ""), configForNeo4j40)
	handleError(err)
	return &NeoDb{d}
}
func CreateSession(d *NeoDb, mode string) neo4j.Session {
	m := neo4j.AccessModeRead
	if mode == "write" {
		m = neo4j.AccessModeWrite
	}
	sessionConfig := neo4j.SessionConfig{AccessMode: m}
	session, err := d.NewSession(sessionConfig)
	handleError(err)
	return session
}

// Here we create a simple function that will take care of errors, helping with some code clean up
func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func DefaultTxMetadata() neo4j.TransactionConfig {
	return neo4j.TransactionConfig{
		Metadata: map[string]interface{}{"user": USER, "datetime": time.Now()},
	}
}