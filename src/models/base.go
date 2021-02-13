package models

//TODO https://www.alexedwards.net/blog/organising-database-access
import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/mindstand/gogm"
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

type Monkey struct {
	gogm.BaseNode
	UUID        string    `gogm:"name=uuid`
	Name        string    `gogm:"name=name"`
	Love        int       `gogm:"name=love"`
	Age         int       `gogm:"name=age"`
	CreatedDate time.Time `gogm:"name=createdDate"`
}

// NewMonkey creates a new instance of monkey in its zero state
func NewMonkey(name string, love int, age int) *Monkey {
	m := &Monkey{
		Name:        name,
		Love:        love,
		Age:         age,
		CreatedDate: time.Now(),
	}
	return m
}

type MonkeyModel struct{}

// Monkeys is a slice of type monkey
type Monkeys []Monkey

// // Persist is a function used to create a new node in neo4j based of a Monkey Struct
// func (mm *monkey) Persist() (string, map[string]interface{}) {

// 	var buffer, returnBuffer bytes.Buffer
// 	buffer.WriteString(fmt.Sprintf("CREATE (n:%s $props)", labelMonkey))
// 	buffer.WriteString(" RETURN ID(n) as id")
// 	props := structs.Map(mm)
// 	keys := make([]string, len(props))
// 	i := 0
// 	for field := range props {
// 		keys[i] = field
// 		i++
// 	}
// 	sort.Strings(keys)
// 	for _, field := range keys {
// 		returnBuffer.WriteString(fmt.Sprintf(", n.%s AS %s", field, field))
// 	}
// 	buffer.Write(returnBuffer.Bytes())
// 	session, err := gogm.NewSession(true)
// 	handleError(err)
// 	defer session.Close()
// 	p := map[string]interface{}{"props": props}
// 	_, err := neo4j.Single(session.Run(buffer.String(), p, neo4j.WithTxMetadata(map[string]interface{}{"user": "neo4j", "datetime": time.Now()})))
// 	if err != nil {
// 		panic(err)
// 	}
// 	return buffer.String(), props
// }

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

func (mm MonkeyModel) GetMonkeys(props []string) []*Monkey {
	session, err := gogm.NewSession(false)
	handleError(err)
	defer session.Close()
	mBuff := getSimpleMatchBuffer("m", labelMonkey)
	// rBuff := getSimpleReturnBuffer("m")
	// rBuff.Write(addFieldsReturn("m", props).Bytes())
	mBuff.WriteString("RETURN m")
	// mBuff.Write(rBuff.Bytes())
	var resultMonkeys []*Monkey
	fmt.Println(mBuff.String())
	err = session.Query(mBuff.String(), nil, &resultMonkeys)
	fmt.Println(resultMonkeys)
	handleError(err)
	return resultMonkeys
}
