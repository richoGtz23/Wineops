package models

//TODO https://www.alexedwards.net/blog/organising-database-access
import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/iancoleman/strcase"
	"github.com/mitchellh/mapstructure"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	uuid "github.com/satori/go.uuid"
)

// // NeoModeler interface used for all the interaction between neo4j and GO
// type NeoModeler interface {
// 	Perist(n *NeoDb) (string, map[string]interface{})
// 	// Update() (string, map[string]interface{})
// 	// Delete() (string, map[string]interface{})
// 	// GetById(uid string) (string, map[string]interface{})
// 	// Filter(map[string]interface{}) (string, map[string]interface{})
// }

const (
	labelMonkey = "Monkey"
)

type monkey struct {
	// gogm.BaseNode
	UID         string    `json:"uid" structs:"uid"`
	Name        string    `json:"name,omitempty" structs:"name,omitempty"`
	Love        int       `json:"love,omitempty" structs:"love,omitempty"`
	Age         int       `json:"age,omitempty" structs:"age,omitempty"`
	CreatedDate time.Time `json:"createdDate" structs:"createdDate,omitnested"`
}

type MonkeyModel struct {
	DB *NeoDb
}

// NewMonkey creates a new instance of monkey in its zero state
func NewMonkey(name string, love int, age int, label string) *monkey {
	m := &monkey{
		UID:         uuid.NewV4().String(),
		Name:        name,
		Love:        love,
		Age:         age,
		CreatedDate: time.Now(),
	}
	return m
}

func (mm MonkeyModel) Create(name string, love int, age int) *monkey {
	m := NewMonkey(name, love, age, "")
	mm.Persist(m)
	return m
}

// Monkeys is a slice of type monkey
type Monkeys []monkey

// Persist is a function used to create a new node in neo4j based of a Monkey Struct
func (mm MonkeyModel) Persist(m *monkey) (string, map[string]interface{}) {

	var buffer, returnBuffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("CREATE (n:%s $props)", labelMonkey))
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
	session := CreateSession(mm.DB, "write")
	defer session.Close()
	p := map[string]interface{}{"props": props}
	_, err := neo4j.Single(session.Run(buffer.String(), p, neo4j.WithTxMetadata(map[string]interface{}{"user": "neo4j", "datetime": time.Now()})))
	if err != nil {
		panic(err)
	}
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

func (mm MonkeyModel) GetMonkeys(props []string) Monkeys {
	session := CreateSession(mm.DB, "read")
	defer session.Close()
	mBuff := getSimpleMatchBuffer("m", labelMonkey)
	rBuff := getSimpleReturnBuffer("m")
	rBuff.Write(addFieldsReturn("m", props).Bytes())
	mBuff.Write(rBuff.Bytes())
	records, err := neo4j.Collect(session.Run(mBuff.String(), nil, neo4j.WithTxMetadata(map[string]interface{}{"user": "neo4j", "datetime": time.Now()})))
	if err != nil {
		panic(err)
	}
	monkeys := make(Monkeys, len(records))
	for i, record := range records {
		m := map[string]interface{}{}
		mk := monkey{}
		for index, key := range record.Keys {
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
