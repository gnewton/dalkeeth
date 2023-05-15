package dalkeeth

import (
	"fmt"
	"log"
	"testing"
)

func Test_SimplePKLookup(t *testing.T) {
	tests := map[int64]bool{
		0:          false,
		VPersonID0: true,
		-1:         false,
	}

	for k, v := range tests {
		log.Println(k, v)
	}
	mdl0, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}
	sess, err := writeTestModelSchema(mdl0)
	if err != nil {
		t.Fatal(err)
	}
	err = writeTestTableRecords(sess)
	if err != nil {
		t.Fatal(err)
	}
	// end setup

	personTbl := sess.TableByKey(TPerson)
	if personTbl == nil {
		t.Fatal("Unable to find table", TPerson)
	}

	idField := personTbl.Field(FId)
	if idField == nil {
		t.Fatal("Unabler to find field with key=", FId)
	}
	nameField := personTbl.Field(FName)
	if nameField == nil {
		t.Fatal(fmt.Errorf("Unabler to find field=%s in table=%s", FName, TPerson))
	}

	rec, err := sess.Get(personTbl, VPersonID0)
	if err != nil {
		t.Fatal(err)
	}

	if rec == nil {
		t.Fatal("Unabler to find record with pid=", VPersonID0)
	}

	t.Log(rec)
	t.Fatal(NotImplemented)
}
