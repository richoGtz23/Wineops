package gql

import (
	"github.com/graphql-go/graphql"
)

// TODO: https://gofiber.io/
type QueryRoot struct {
	Query *graphql.Object
}

func NewQueryRoot() *QueryRoot {
	resolver := Resolver{}

	qRoot := QueryRoot{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"Monkeys": &graphql.Field{
					Type:        graphql.NewList(MonkeyType),
					Description: "List of monkeys",
					Resolve:     resolver.MonkeysResolver,
				},
			},
		}),
	}
	return &qRoot
}
