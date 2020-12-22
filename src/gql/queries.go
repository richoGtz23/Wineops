package gql

import (
	"github.com/RichoGtz23/Wineops/src/models"
	"github.com/graphql-go/graphql"
)

type QueryRoot struct {
	Query *graphql.Object
}

func NewQueryRoot(db *models.NeoDb) *QueryRoot {
	resolver := Resolver{db: db}

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
