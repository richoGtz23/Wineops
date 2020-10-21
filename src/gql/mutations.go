package gql

import (
	"github.com/RichoGtz23/Wineops/src/models"
	"github.com/graphql-go/graphql"
)

type MutationRoot struct {
	Mutation *graphql.Object
}

func NewMutationRoot(db *models.NeoDb) *MutationRoot {
	resolver := Resolver{db: db}

	m_root := MutationRoot{
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				"CreateMonkey": &graphql.Field{
					Type:        MonkeyType,
					Description: "Creates a new monkey",
					Args: graphql.FieldConfigArgument{
						"name": &graphql.ArgumentConfig{
							Type:        graphql.NewNonNull(graphql.String),
							Description: "The name of the monkey",
						},
						"age": &graphql.ArgumentConfig{
							Type:        graphql.NewNonNull(graphql.Int),
							Description: "Age of the monkey",
						},
						"love": &graphql.ArgumentConfig{
							Type:        graphql.Int,
							Description: "Love level of the monkey",
						},
					},
					Resolve: resolver.CreateMonkeyResolver,
				},
			},
		}),
	}
	return &m_root
}
