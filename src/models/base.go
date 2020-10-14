package models

//TODO https://www.alexedwards.net/blog/organising-database-access
import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/RichoGtz23/Wineops/src/neodb"
	"github.com/fatih/structs"
	"github.com/iancoleman/strcase"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	uuid "github.com/satori/go.uuid"
)

type baseModel struct {
	ID    int64 `json:"id,omitempty" structs:"id"`
	label string
}

// NeoModeler interface used for all the interaction between neo4j and GO
type NeoModeler interface {
	setLabel(label string)
	Create() (string, map[string]interface{})
	// Update() (string, map[string]interface{})
	// Delete() (string, map[string]interface{})
	// GetById(uid string) (string, map[string]interface{})
	// Filter(map[string]interface{}) (string, map[string]interface{})
}

func ExtractLabel(n NeoModeler) (interface{}, string) {
	v := reflect.ValueOf(n)
	if reflect.TypeOf(n).Kind() == reflect.Ptr {
		v = v.Elem()
	}
	t := v.Type()
	l := strings.Split(t.String(), ".")
	return t, strcase.ToCamel(l[len(l)-1])
}

type monkey struct {
	baseModel   `structs:",omitnested"`
	UID         string    `json:"uid" structs:"uid"`
	Name        string    `json:"name,omitempty" structs:"name,omitempty"`
	Love        int       `json:"love,omitempty" structs:"love,omitempty"`
	Age         int       `json:"age,omitempty" structs:"age,omitempty"`
	CreatedDate time.Time `json:"createdDate" structs:"createdDate,omitnested"`
}

// NewMonkey creates a new instance of monkey in its zero state
func NewMonkey(name string, love int, age int, label string) *monkey {
	m := &monkey{
		baseModel:   baseModel{},
		UID:         uuid.NewV4().String(),
		Name:        name,
		Love:        love,
		Age:         age,
		CreatedDate: time.Now(),
	}
	m.setLabel(label)
	return m
}

// Monkeys is a slice of type monkey
type Monkeys []monkey

func (m *monkey) setLabel(label string) {
	if m.label != "" {
		return
	}
	if label == "" {
		l := strings.Split(reflect.TypeOf(m).String(), ".")
		label := strcase.ToCamel(l[len(l)-1])
		m.label = label
	} else {
		m.label = label
	}
}

func getType(t interface{}) string {
	return reflect.TypeOf(t).Name()
}

// CreateNodeCypher is a function used to create a new node in neo4j based of a Monkey Struct
func (m *monkey) Create() (string, map[string]interface{}) {

	con := neodb.CreateConnection()
	defer con.Close()

	var buffer, returnBuffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("CREATE (n:%s $props)", m.label))
	buffer.WriteString(" RETURN ID(n) as id")
	props := structs.Map(m)
	keys := make([]string, len(props))
	i := 0
	for field := range props {
		keys[i] = field
		i++
	}
	sort.Strings(keys)
	for _, field := range keys {
		returnBuffer.WriteString(fmt.Sprintf(", n.%s AS %s", field, field))
	}
	buffer.Write(returnBuffer.Bytes())
	session := neodb.CreateSession(con, "write")
	defer session.Close()
	p := map[string]interface{}{"props": props}
	r, _ := neo4j.Single(session.Run(buffer.String(), p, neo4j.WithTxMetadata(map[string]interface{}{"user": neodb.USER, "datetime": time.Now()})))
	// for index, key := range r.Keys() {
	// 	if index == 0 {
	// 		continue
	// 	}
	// 	fmt.Println(key, r.GetByIndex(index), props[key])
	// }
	m.ID = r.GetByIndex(0).(int64)
	return buffer.String(), props
}

func getSimpleReturnBuffer(nName string) *bytes.Buffer {
	return bytes.NewBuffer([]byte(fmt.Sprintf("RETURN ID(%s) as id", nName)))
}

func getSimpleMatchBuffer(nName string, label string) *bytes.Buffer {
	return bytes.NewBuffer([]byte(fmt.Sprintf("MATCH (%s :%s) ", nName, strcase.ToCamel(strings.ToLower(label)))))
}
func addWhere(b *bytes.Buffer) {
	b.WriteString("WHERE ")
}
func addFieldsReturn(nName string, props []string) *bytes.Buffer {
	b := bytes.NewBuffer(nil)
	for _, fld := range props {
		b.WriteString(fmt.Sprintf(", %s.%s as %s", nName, strcase.ToLowerCamel(fld), strcase.ToCamel(fld)))
	}
	return b
}

func GetMonkeys(props []string) Monkeys {
	con := neodb.CreateConnection()
	defer con.Close()
	session := neodb.CreateSession(con, "read")
	defer session.Close()
	mBuff := getSimpleMatchBuffer("m", getType(monkey{}))
	rBuff := getSimpleReturnBuffer("m")
	rBuff.Write(addFieldsReturn("m", props).Bytes())
	mBuff.Write(rBuff.Bytes())
	records, err := neo4j.Collect(session.Run(mBuff.String(), nil, neo4j.WithTxMetadata(map[string]interface{}{"user": neodb.USER, "datetime": time.Now()})))
	if err != nil {
		panic(err)
	}
	monkeys := make(Monkeys, len(records))
	for i, record := range records {
		m := map[string]interface{}{}
		mk := monkey{}
		for index, key := range record.Keys() {
			m[key] = record.GetByIndex(index)
		}
		err := mapstructure.Decode(m, &mk)
		if err != nil {
			panic(err)
		}
		monkeys[i] = mk

	}
	return monkeys
}
