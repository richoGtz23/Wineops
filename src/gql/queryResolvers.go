package gql

import (
	"fmt"

	. "github.com/RichoGtz23/Wineops/src/models"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

type Resolver struct{}

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

func (r *Resolver) MonkeysResolver(p graphql.ResolveParams) (interface{}, error) {
	fields, _ := getSelectedFields([]string{"Monkeys"}, p)
	f := make([]string, len(fields["Monkeys"]))
	i := 0
	for k, v := range fields["Monkeys"] {
		if v == true {
			f[i] = k
			i++
		}
	}
	return MonkeyModel{}.GetMonkeys(f), nil
}
