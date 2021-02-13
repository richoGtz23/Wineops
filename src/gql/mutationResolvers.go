package gql

import (
	. "github.com/RichoGtz23/Wineops/src/models"
	"github.com/graphql-go/graphql"
	"github.com/mindstand/gogm"
)

func (r *Resolver) CreateMonkeyResolver(params graphql.ResolveParams) (interface{}, error) {
	monkey := NewMonkey(params.Args["name"].(string), params.Args["love"].(int), params.Args["age"].(int))
	session, err := gogm.NewSession(false)
	if err != nil {
		panic(err)
	}
	defer session.Close()
	session.Save(monkey)
	return monkey, nil
}
