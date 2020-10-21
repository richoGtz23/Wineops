package gql

import (
	"github.com/graphql-go/graphql"
)

func (r *Resolver) CreateMonkeyResolver(params graphql.ResolveParams) (interface{}, error) {
	// newMonkey := models.NewMonkey(params.Args["name"].(string),
	// 	params.Args["love"].(int), params.Args["age"].(int), "")
	// newMonkey.Create()
	monkey := r.db.CreateMonkey(params.Args["name"].(string), params.Args["love"].(int), params.Args["age"].(int))
	return monkey, nil
}
