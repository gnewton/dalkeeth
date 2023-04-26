package dalkeeth

import (
	"testing"
)

func Test00(t *testing.T) {
	setupTest()
	model, err := defineTestModel()
	if err != nil {
		t.Fatal(err)
	}

	//age, ok := model.fieldTableMap["persons.age"]
	age, ok := model.fieldTableMap["addresses.street"]
	if !ok {
		t.Log(model.fieldTableMap)
		t.Fatal("Unable to find persons.age field")
	}

	name := Field{name: "name", fieldType: StringType}

	q := SelectQuery{
		Fields:         []AField{age, &name},
		Pks:            []int64{54, 767},
		Where:          WN(&name, IsNotNull),
		GroupBy:        []*Field{&name},
		Having:         W(age, GT, 100),
		Offset:         1200,
		Limit:          100,
		GlobalOrdering: ASC,
		OrderByFields:  []*FieldOrdered{},
	}

	if err := q.Validate(model); err != nil {
		//t.Error(err) FIXX
	}

}
