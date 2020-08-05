// TODO: Add https: https://medium.com/rungo/secure-https-servers-in-go-a783008b36da
// TODO: Add http 2: https://marcofranssen.nl/build-a-go-webserver-on-http-2-using-letsencrypt/
// TODO: compose https://medium.com/@magbicaleman/go-graphql-and-neo4j-6d65b28736cd
// TODO: auth graphql https://www.howtographql.com/graphql-go/6-authentication/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"

	neodb "github.com/RichoGtz23/Wineops/src/utils"

	"github.com/iancoleman/strcase"
	"github.com/julienschmidt/httprouter"
	UUID "github.com/satori/go.uuid"
)

// NeoObject type are structs that are compatible to use predefined queries for CUDR on Neo4j
type NeoObject interface {
	createNodeCypher()
}

// NeoBase Contains essential functionality to create Struct and be able to save and load from neo4j database
type NeoBase struct {
	Uid string `json:"uid"`
}

// Monkey type is amazing
type Monkey struct {
	NeoBase
	Name string `json:"name"`
	Love int    `json:"love"`
	Age  int    `json:"age"`
}

func (m *Monkey) createNodeCypher() (string, map[string]interface{}) {

	con := neodb.CreateConnection()
	defer con.Close()
	labels := strings.Split(reflect.TypeOf(m).String(), ".")
	label := strcase.ToCamel(labels[len(labels)-1])

	var buffer, returnBuffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("CREATE (n:%s $props)", label))
	buffer.WriteString(" RETURN ID(n)")
	props := make(map[string]interface{})
	reflection := reflect.ValueOf(m).Elem()
	var varName string
	var varValue interface{}
	for i := 0; i < reflection.NumField(); i++ {
		if reflection.Field(i).Kind() == reflect.Struct {
			valueValue := reflect.ValueOf(reflection.Field(i).Interface())
			varName = valueValue.Type().Field(0).Name
			varValue = valueValue.Field(0).Interface()
		} else {
			varName = reflection.Type().Field(i).Name
			varValue = reflection.Field(i).Interface()
		}
		varName = strcase.ToLowerCamel(varName)
		returnBuffer.WriteString(fmt.Sprintf(", n.%s AS %s", varName, varName))
		props[varName] = varValue
	}
	buffer.Write(returnBuffer.Bytes())

	st := neodb.PrepareSatement(buffer.String(), con)
	neodb.ExecuteStatement(st, props)
	return buffer.String(), props
}

func indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	u := UUID.NewV4()
	u2 := UUID.NewV4()
	uu := u.String()
	uu2 := u2.String()
	m := Monkey{
		NeoBase: NeoBase{Uid: uu},
		Name:    "Karinepe",
		Love:    100,
		Age:     24,
	}
	m2 := Monkey{
		NeoBase: NeoBase{Uid: uu2},
		Name:    "Richo",
		Love:    100,
		Age:     25,
	}
	cypher, props := m.createNodeCypher()
	b, err := json.Marshal(props)
	fmt.Fprintf(w, cypher)
	fmt.Fprintf(w, string(b))
	monkeys := []Monkey{m, m2}
	b, err = json.Marshal(monkeys)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(w, string(b))
}

func main() {
	router := httprouter.New()
	router.GET("/", indexHandler)
	// print env
	env := os.Getenv("APP_ENV")
	if env == "production" {
		log.Println("Running api server in production mode")
	} else {
		log.Println("Running api server in dev mode")
	}

	log.Fatal(http.ListenAndServe(":8080", router))
}
