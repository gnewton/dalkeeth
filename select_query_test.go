package dalkeeth

import (
	"testing"
)

func TestSelectQuery_FullStruct(t *testing.T) {
	setupTest()
	model, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}

	ageField := model.TableField(TAddress, FStreet)
	if ageField == nil {
		t.Log(model.fieldTableMap)
		t.Fatal("Unable to find persons.age field")
	}

	addressTable, ok := model.tablesMap[TAddress]
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
