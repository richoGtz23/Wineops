// TODO: Add https: https://medium.com/rungo/secure-https-servers-in-go-a783008b36da
// TODO: Add http 2: https://marcofranssen.nl/build-a-go-webserver-on-http-2-using-letsencrypt/
// TODO: compose https://medium.com/@magbicaleman/go-graphql-and-neo4j-6d65b28736cd
// TODO: auth graphql https://www.howtographql.com/graphql-go/6-authentication/
// TODO: clean devops Dockerfile for go and docker-compose.yml
// TODO: clean base.go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

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
	file, err := os.Open("/run/secrets/APP_ENV")
	if err != nil {
		log.Println("No ENV is set")
		return
	}
	envData := map[string]string{}
	bytes, err := ioutil.ReadAll(file)
	err = json.Unmarshal(bytes, &envData)
	if err != nil {
		log.Println("Couldn't read environment data")
		return
	}
	env := envData["APP_ENV"]
	if env == "production" {
		log.Println("Running api server in production mode")
	} else {
		log.Println("Running api server in dev mode")
	}
	neo4jHost := envData["NEO4J_HOST"]
	neo4jPassword := envData["NEO4J_PASSWORD"]
	neo4jUser := envData["NEO4J_USER"]
	neo4jPort := envData["NEO4J_PORT"]
	if neo4jHost == "" || neo4jPassword == "" || neo4jUser == "" || neo4jPort == "" {
		log.Println("No neo4j ENV variables defined")
		return
	}

	router := initializeAPI(neo4jHost, neo4jPassword, neo4jUser, neo4jPort)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func initializeAPI(neo4jHost, neo4jPassword, neo4jUser, neo4jPort string) *chi.Mux {
	router := chi.NewRouter()
	intPort, err := strconv.Atoi(neo4jPort)
	if err != nil {
		panic(err)
	}
	models.CreateConnection(neo4jHost, neo4jPassword, neo4jUser, intPort)
	rootQuery := gql.NewQueryRoot()
	rootMutation := gql.NewMutationRoot()
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

	// Create the graphql route with a Server method to handle it
	router.Post("/graphql", s.GraphQL())
	// router.Get("/favicon", http.NotFoundHandler())
	router.Get("/graphiql", graphiql.ServeGraphiQL)

	return router
}
