// TODO: Add https: https://medium.com/rungo/secure-https-servers-in-go-a783008b36da
// TODO: Add http 2: https://marcofranssen.nl/build-a-go-webserver-on-http-2-using-letsencrypt/
// TODO: compose https://medium.com/@magbicaleman/go-graphql-and-neo4j-6d65b28736cd
// TODO: auth graphql https://www.howtographql.com/graphql-go/6-authentication/
// TODO: clean devops Dockerfile for go and docker-compose.yml
// TODO: clean base.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/RichoGtz23/Wineops/src/gql"
	"github.com/RichoGtz23/Wineops/src/models"
	"github.com/RichoGtz23/Wineops/src/server"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/graphql-go/graphql"
	"github.com/mnmtanish/go-graphiql"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "production" {
		log.Println("Running api server in production mode")
	} else {
		log.Println("Running api server in dev mode")
	}

	router, db := initializeAPI()
	defer db.Close()

	log.Fatal(http.ListenAndServe(":8080", router))
}

func initializeAPI() (*chi.Mux, *models.NeoDb) {
	// Create a new router
	router := chi.NewRouter()

	// Create a new connection to our pg database
	db := models.CreateConnection()

	// Create our root query for graphql
	rootQuery := gql.NewQueryRoot(db)
	rootMutation := gql.NewMutationRoot(db)
	// Create a new graphql schema, passing in the the root query
	sc, err := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    rootQuery.Query,
			Mutation: rootMutation.Mutation,
		},
	)
	if err != nil {
		fmt.Println("Error creating schema: ", err)
	}

	// Create a server struct that holds a pointer to our database as well
	// as the address of our graphql schema
	s := server.Server{
		GqlSchema: &sc,
	}

	// Add some middleware to our router
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // set content-type headers as application/json
		middleware.Logger,          // log api request calls
		middleware.DefaultCompress, // compress results, mostly gzipping assets and json
		middleware.StripSlashes,    // match paths with a trailing slash, strip it, and continue routing through the mux
		middleware.Recoverer,       // recover from panics without crashing server
	)

	// var c = cors.New(cors.Options{
	// 	AllowedOrigins: []string{"*"},
	// })
	// var h = handler.New(&handler.Config{
	// 	Schema: &Schema,
	// 	Pretty: true,
	// })
	// serveMux := http.NewServeMux()
	// serveMux.Handle("/favicon", http.NotFoundHandler())
	// serveMux.Handle("/graphql", c.Handler(h))
	// serveMux.HandleFunc("/graphiql", graphiql.ServeGraphiQL)

	// Create the graphql route with a Server method to handle it
	router.Post("/graphql", s.GraphQL())
	// router.Get("/favicon", http.NotFoundHandler())
	router.Get("/graphiql", graphiql.ServeGraphiQL)

	return router, db
}
