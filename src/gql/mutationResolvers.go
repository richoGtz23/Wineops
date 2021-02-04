package gql

import (
	. "github.com/RichoGtz23/Wineops/src/models"
	"github.com/graphql-go/graphql"
)

func (r *Resolver) CreateMonkeyResolver(params graphql.ResolveParams) (interface{}, error) {
	monkeyModel := MonkeyModel{DB: r.db}
	monkey := monkeyModel.Create(params.Args["name"].(string), params.Args["love"].(int), params.Args["age"].(int))
	return monkey, nil
}
