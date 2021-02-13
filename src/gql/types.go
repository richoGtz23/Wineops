package gql

import "github.com/graphql-go/graphql"

// MonkeyType defines the exposed object for a monkey in GraphQL
var MonkeyType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Monkey",
		Fields: graphql.Fields{
			"uuid": &graphql.Field{
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
