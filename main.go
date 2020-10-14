// TODO: Add https: https://medium.com/rungo/secure-https-servers-in-go-a783008b36da
// TODO: Add http 2: https://marcofranssen.nl/build-a-go-webserver-on-http-2-using-letsencrypt/
// TODO: compose https://medium.com/@magbicaleman/go-graphql-and-neo4j-6d65b28736cd
// TODO: auth graphql https://www.howtographql.com/graphql-go/6-authentication/
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RichoGtz23/Wineops/src/models"
	"github.com/RichoGtz23/Wineops/src/neodb"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/handler"
	"github.com/mnmtanish/go-graphiql"
	"github.com/rs/cors"
)

func getSelectedFields(selected []string, params graphql.ResolveParams) (map[string]map[string]interface{}, error) {
	fieldASTs := params.Info.FieldASTs
	if len(fieldASTs) == 0 {
		return nil, fmt.Errorf("getSelectedFields: ResolveParams has no fields")
	}
	selections := map[string]map[string]interface{}{}
	for _, s := range selected {
		for _, field := range fieldASTs {
			if field.Name.Value == s {
				var err error
				selections[s], err = selectedFieldsFromSelections(params, field.SelectionSet.Selections)
				if err != nil {
					return nil, err
				}
				break
			}
		}
	}
	return selections, nil
}

func selectedFieldsFromSelections(params graphql.ResolveParams, selections []ast.Selection) (selected map[string]interface{}, err error) {
	selected = map[string]interface{}{}

	for _, s := range selections {
		switch s := s.(type) {
		case *ast.Field:
			if s.SelectionSet == nil {
				selected[s.Name.Value] = true
			} else {
				selected[s.Name.Value], err = selectedFieldsFromSelections(params, s.SelectionSet.Selections)
				if err != nil {
					return
				}
			}
		case *ast.FragmentSpread:
			n := s.Name.Value
			frag, ok := params.Info.Fragments[n]
			if !ok {
				err = fmt.Errorf("getSelectedFields: no fragment found with name %v", n)

				return
			}
			selected[s.Name.Value], err = selectedFieldsFromSelections(params, frag.GetSelectionSet().Selections)
			if err != nil {
				return
			}
		default:
			err = fmt.Errorf("getSelectedFields: found unexpected selection type %v", s)

			return
		}
	}

	return
}

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
var mutationType = graphql.NewObject(graphql.ObjectConfig{
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
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				newMonkey := models.NewMonkey(params.Args["name"].(string),
					params.Args["love"].(int), params.Args["age"].(int), "")
				newMonkey.Create()
				// fmt.Println(models.ExtractLabel(newMonkey))
				return newMonkey, nil
			},
		},
	},
})

var queryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.Fields{
		"Monkeys": &graphql.Field{
			Type:        graphql.NewList(MonkeyType),
			Description: "List of monkeys",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				fields, _ := getSelectedFields([]string{"Monkeys"}, p)
				f := make([]string, len(fields["Monkeys"]))
				i := 0
				for k, v := range fields["Monkeys"] {
					if v == true {
						f[i] = k
						i++
					}
				}
				return models.GetMonkeys(f), nil
			},
		},
	},
})
var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    queryType,
	Mutation: mutationType,
})

var c = cors.New(cors.Options{
	AllowedOrigins: []string{"*"},
})
var h = handler.New(&handler.Config{
	Schema: &Schema,
	Pretty: true,
})

func main() {
	con := neodb.CreateConnection()
	defer con.Close()

	serveMux := http.NewServeMux()
	env := os.Getenv("APP_ENV")
	if env == "production" {
		log.Println("Running api server in production mode")
	} else {
		log.Println("Running api server in dev mode")
	}
	serveMux.Handle("/favicon", http.NotFoundHandler())
	serveMux.Handle("/graphql", c.Handler(h))
	serveMux.HandleFunc("/graphiql", graphiql.ServeGraphiQL)

	log.Fatal(http.ListenAndServe(":8080", serveMux))
}
