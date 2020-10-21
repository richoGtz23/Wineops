package gql

import "github.com/graphql-go/graphql"

// GQL
var MonkeyType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Monkey",
		Fields: graphql.Fields{
			"uid": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"love": &graphql.Field{
				Type: graphql.Int,
			},
			"age": &graphql.Field{
				Type: graphql.Int,
			},
			"createdDate": &graphql.Field{
				Type: graphql.DateTime,
			},
		},
	},
)
