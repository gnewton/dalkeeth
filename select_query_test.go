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
	//age, ok := model.fieldTableMap["addresses.street"]
	ageField, ok := model.fieldTableMap[TAddress+"."+FStreet]
	if !ok {
		t.Log(model.fieldTableMap)
		t.Fatal("Unable to find persons.age field")
	}

	addressTable, ok := model.tablesMap[TAddressK]
	if !ok {
		t.Fatal("Unable to find addresses table")
	}

	nameField := Field{name: "name", fieldType: StringType}
	// Assumed to be from

	q := SelectQuery{
		Fields:         []AField{ageField, &nameField},
		From:           []*Table{addressTable},
		Pks:            []int64{54, 767},
		Where:          WN(&nameField, IsNotNull),
		GroupBy:        []*Field{&nameField},
		Having:         W(ageField, GT, 100),
		Offset:         1200,
		Limit:          100,
		GlobalOrdering: ASC,
		OrderByFields:  []*FieldOrdered{},
	}

	if err := q.Validate(model); err != nil {
		//t.Error(err) FIXX
	}

}
