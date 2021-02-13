// TODO: official NEO4j https://medium.com/@angadsharma1016/optimizing-go-neo4j-concurrency-patterns-810dff25f88f
package models

import (
	"time"

	"github.com/mindstand/gogm"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

const POOLSIZE = 50

func CreateConnection(neo4jHost, neo4jPassword, neo4jUser string, neo4jPort int) {
	config := gogm.Config{
		IndexStrategy: gogm.VALIDATE_INDEX, //other options are ASSERT_INDEX and IGNORE_INDEX
		PoolSize:      POOLSIZE,
		Port:          neo4jPort,
		IsCluster:     false, //tells it whether or not to use `bolt+routing`
		Host:          neo4jHost,
		Password:      neo4jPassword,
		Username:      neo4jUser,
	}
	err := gogm.Init(&config, &Monkey{})
	if err != nil {
		panic(err)
	}
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
