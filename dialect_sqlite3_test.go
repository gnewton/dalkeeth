package dalkeeth

import (
	"testing"
)

func Test_DeleteSql_Raw(t *testing.T) {
	sess, err := initAndWriteTestTableSchema()
	if err != nil {
		t.Fatal(err)
	}
	err = writeTestTableRecords(sess)
	if err != nil {
		t.Fatal(err)
	}

	personTbl := sess.TableByKey(TPerson)
	if personTbl == nil {
		t.Fatal("Unable to find table")
	}

	idField := personTbl.Field(FId)
	if idField == nil {
		t.Fatal("Unabler to find field with key=", FId)
	}
	nameField := personTbl.Field(FName)
	if nameField == nil {
		t.Fatal("Unabler to find field with key=", FName)
	}

	r := personTbl.NewRecord()
	//r.SetValues([]*Value{})
	r.SetValues([]*Value{{field: idField, value: 43}, {field: nameField, value: "Bob"}})

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
