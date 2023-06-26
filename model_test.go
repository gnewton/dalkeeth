package dalkeeth

import (
	"fmt"
	"testing"
)

func Test_Model_AddForeignKey_UnknownForeignKeyField(t *testing.T) {
	setupTest()
	model, persons, addresses, err := addForeignKey_Setup()

	if err != nil {
		t.Error(err)
	}

	if model.AddForeignKey(persons, "foo", addresses, FId) == nil {
		t.Fatal(fmt.Errorf("Failed identifying incorrect field"))
	}
}

func Test_Model_AddForeignKey_UnknownForeignKeyFieldOtherField(t *testing.T) {
	setupTest()
	model, err := testModel0()

	if err != nil {
		t.Error(err)
	}

	persons := model.TableByKey(TPerson)

	if persons == nil {
		t.Log(model.tablesMap)
		t.Fatal(fmt.Errorf("Table key %s not found by manager but should be found", TPerson))
	}

	addresses := model.TableByKey(TAddress)

	if addresses == nil {
		t.Fatal(fmt.Errorf("Table key %s not found by manager but should be found", TAddress))
	}

	if model.AddForeignKey(persons, FId, addresses, "foo") == nil {
		t.Fatal(fmt.Errorf("Failed identifying incorrect field"))
	}
}
