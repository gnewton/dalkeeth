package dalkeeth

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"testing"
)

var ShouldHaveFailed = errors.New("Should have failed.")

func TestSession_Table_EmptyString(t *testing.T) {
	setupTest()
	model, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}
	// end setup

	tName := ""
	tbl := model.TableByKey(tName)
	if tbl != nil {
		t.Fatal(fmt.Errorf("Table with empty string key"))
	}

}

func TestSession_Table_UnknownString(t *testing.T) {
	setupTest()
	sess, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}
	// end setup

	tName := "not-existing-table-foo"
	tbl := sess.TableByKey(tName)
	if tbl != nil {
		t.Fatal(fmt.Errorf("Table should not exist: %s", tName))
	}

}

func TestSession_AddTableToFrozenModel(t *testing.T) {
	setupTest()
	model, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}

	// end setup

	_, err = model.NewTable(TPerson)
	if err == nil {
		t.Fatal(errors.New("Should fail: cannot add table to frozen model"))
	}

}

func TestSession_CreateTablesSQL_NilDialect(t *testing.T) {
	setupTest()
	model, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}

	// end setup
	sess, err := NewSession(model)
	if err != nil {
		t.Fatal(err)
	}
	sess.dialect = nil

	_, err = sess.createTablesSQL()
	if err == nil {
		t.Fatal(ShouldHaveFailed)
	}
}

func TestSession_CreateTableIndexesSQL_NilDialect(t *testing.T) {
	setupTest()
	model, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}

	// end setup
	sess, err := NewSession(model)
	if err != nil {
		t.Fatal(err)
	}
	sess.dialect = nil

	_, err = sess.createTableIndexesSQL()
	if err == nil {
		t.Fatal(ShouldHaveFailed)
	}
}

//func TestSession_AddForeignKey_NilTable(t *testing.T) {
//	setupTest()
//	t.Fatal(NotImplemented)
//}

func Test_SessionInitTables2(t *testing.T) {
	setupTest()
	_, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Session_Session_SaveTx(t *testing.T) {
	setupTest()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	mdl0, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}

	sess, err := writeTestModelSchema(mdl0)
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()

	persons := sess.TableByKey(TPerson)
	if persons == nil {
		t.Error(errors.New("Persons cannot be found: is nil"))
	}

	sess.dialect = new(DialectSqlite3)

	pk := int64(23)
	rec := persons.NewRecord()
	err = rec.SetValue(FId, pk)
	if err != nil {
		t.Fatal(err)
	}
	err = rec.SetValue(FName, "Fred")
	if err != nil {
		t.Fatal(err)
	}

	err = rec.SetValue(FAge, 54)
	if err != nil {
		t.Fatal(err)
	}

	err = sess.Begin()
	if err != nil {
		t.Fatal(err)
	}

	err = sess.SaveTx(rec)
	if err != nil {
		t.Fatal(err)
	}

	err = sess.Commit()
	if err != nil {
		t.Fatal(err)
	}

	valid, err := recordExists(sess.db, rec.table.name, pk)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal(errors.New("Value not in database"))
	}

}

func Test_Session_Session_Save(t *testing.T) {
	setupTest()
	mdl0, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}
	sess, err := writeTestModelSchema(mdl0)
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()

	persons := sess.TableByKey(TPerson)
	if persons == nil {
		t.Fatal(fmt.Errorf("Table key %s not found by manager but should be found", TPerson))
	}

	sess.dialect = new(DialectSqlite3)

	pk := int64(23)
	rec := persons.NewRecord()
	err = rec.SetValue(FId, pk)
	if err != nil {
		t.Fatal(err)
	}
	err = rec.SetValue(FName, "Fred")
	if err != nil {
		t.Fatal(err)
	}

	err = rec.SetValue(FAge, 54)
	if err != nil {
		t.Fatal(err)
	}

	err = sess.Save(rec)
	if err != nil {
		t.Fatal(err)
	}

	// See if we can read the record that was just written
	valid, err := recordExists(sess.db, rec.table.name, pk)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal(errors.New("Value not in database"))
	}
}

func Test_Session_Session_Batch(t *testing.T) {
	setupTest()
	mdl0, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}
	sess, err := writeTestModelSchema(mdl0)
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()

	persons := sess.TableByKey(TPerson)
	if persons == nil {
		t.Fatal(fmt.Errorf("Table key %s not found by manager but should be found", TPerson))
	}

	sess.dialect = new(DialectSqlite3)

	records, err := twoPersonRecords(persons)
	if err != nil {
		t.Fatal(err)
	}

	err = sess.Batch(records)
	if err != nil {
		t.Fatal(err)
	}

	// See if the 2 records added are readable
	valid, err := recordExists(sess.db, records[0].table.name, VPersonID0)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal(errors.New("Value not in database"))
	}
	valid, err = recordExists(sess.db, records[1].table.name, VPersonID1)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal(errors.New("Value not in database"))
	}
}

func Test_Session_Session_Begin_DBNil(t *testing.T) {
	setupTest()
	mdl0, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}
	sess, err := writeTestModelSchema(mdl0)
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()
	sess.db = nil
	err = sess.Begin()
	if err == nil {
		t.Fatal(ShouldHaveFailed)
	}
}

func Test_Session_Session_Begin_DoubleTx(t *testing.T) {
	setupTest()
	mdl0, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}
	sess, err := writeTestModelSchema(mdl0)
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()

	err = sess.Begin()
	if err != nil {
		t.Fatal(err)
	}

	// Should fail
	err = sess.Begin()
	if err == nil {
		t.Fatal(ShouldHaveFailed)
	}
	log.Println(err)
}

func Test_Session_Session_Commit_NilTx(t *testing.T) {
	setupTest()
	mdl0, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}
	sess, err := writeTestModelSchema(mdl0)
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()

	// Should fail
	err = sess.Commit()
	if err == nil {
		t.Fatal(ShouldHaveFailed)
	}
	log.Println(err)
}

func Test_Session_Session_BatchMany(t *testing.T) {
	setupTest()
	mdl0, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}
	sess, err := writeTestModelSchema(mdl0)
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()

	persons := sess.TableByKey(TPerson)
	if persons == nil {
		t.Fatal(fmt.Errorf("Table key %s not found by manager but should be found", TPerson))
	}

	records, err := nPersonRecords(persons, 10000)
	if err != nil {
		t.Fatal(err)
	}

	err = sess.Batch(records)
	if err != nil {
		t.Fatal(err)
	}

	valid, err := recordExists(sess.db, records[0].table.name, VPersonID0)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal(errors.New("Value not in database"))
	}
	log.Println("Found", VPersonID0, "in DB")

	valid, err = recordExists(sess.db, records[1].table.name, VPersonID1)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal(errors.New("Value not in database"))
	}
	log.Println("Found", VPersonID1, "in DB")
}

func Test_Session_Session_Save_MissingNotNullValue(t *testing.T) {
	setupTest()
	db, err := openTestDB()
	if err != nil {
		t.Error(err)
	}

	model, persons, _, err := addForeignKey_Setup()

	if err != nil {
		t.Error(err)
	}
	sess, err := NewSession(model)
	if err != nil {
		t.Error(err)
	}
	sess.db = db
	sess.dialect = new(DialectSqlite3)

	rec := persons.NewRecord()
	err = rec.SetValue(FId, 23)
	if err != nil {
		t.Fatal(err)
	}
	err = rec.SetValue(FAge, 54)
	if err != nil {
		t.Fatal(err)
	}

	err = sess.Begin()
	if err != nil {
		t.Fatal(err)
	}

	err = sess.Save(rec)
	if err == nil {
		t.Fatal(ShouldHaveFailed)
	}
}

func Test_Session_Session_Save_MissingPK(t *testing.T) {
	setupTest()
	db, err := openTestDB()
	if err != nil {
		t.Error(err)
	}
	model, persons, _, err := addForeignKey_Setup()

	if err != nil {
		t.Error(err)
	}

	sess, err := NewSession(model)
	if err != nil {
		t.Error(err)
	}
	sess.db = db
	sess.dialect = new(DialectSqlite3)

	rec := persons.NewRecord()

	err = rec.SetValue(FName, "Fred")
	if err != nil {
		t.Fatal(err)
	}

	err = rec.SetValue(FAge, 54)
	if err != nil {
		t.Fatal(err)
	}

	err = sess.Begin()
	if err != nil {
		t.Fatal(err)
	}

	err = sess.Save(rec)
	if err == nil {
		t.Fatal(ShouldHaveFailed)
	}
}

func Test_Session_Session_Get(t *testing.T) {
	setupTest()
	mdl0, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}
	sess, err := writeTestModelSchema(mdl0)
	if err != nil {
		t.Fatal(err)
	}
	defer sess.Close()

	persons := sess.TableByKey(TPerson)
	if persons == nil {
		t.Fatal(fmt.Errorf("Table key %s not found by manager but should be found", TPerson))
	}

	sess.dialect = new(DialectSqlite3)

	records, err := twoPersonRecords(persons)
	if err != nil {
		t.Fatal(err)
	}

	err = sess.Batch(records)
	if err != nil {
		t.Fatal(err)
	}
	valid, err := recordExists(sess.db, records[0].table.name, VPersonID0)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal(errors.New("Value not in database"))
	}
	log.Println("Found", VPersonID0, "in DB")

	valid, err = recordExists(sess.db, records[1].table.name, VPersonID1)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal(errors.New("Value not in database"))
	}
	log.Println("Found", VPersonID1, "in DB")

	// end setup

	rec, err := sess.Get(persons, VPersonID0)
	if err != nil {
		t.Fatal(err)
	}

	log.Println("VALUE=", rec.values[0].value)
	log.Println("VALUE=", rec.values[0].value)

	switch v := rec.values[0].value.(type) {
	case *int64:
		t.Log("is int", v)
	}
	//t.Log(rec.values[0].value)
	//t.Log(rec.values[1].value)
	//t.Log(rec.values[2].value)
}

func TestSession_InstantiateModel_ReadOnly(t *testing.T) {
	setupTest()
	model, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}
	sess, err := NewSession(model)
	if err != nil {
		t.Fatal(err)
	}

	// end setup

	err = sess.WriteModelTableSchemaToDB()
	if err == nil {
		t.Fatal(ShouldHaveFailed)
	}
	//t.Log(err)
}

func openTestDB() (*sql.DB, error) {
	return sql.Open("sqlite3", ":memory:")
}
func addForeignKey_Setup() (*Model, *Table, *Table, error) {
	model, err := testModel0()

	if err != nil {
		return nil, nil, nil, err
	}

	persons := model.TableByKey(TPerson)

	if persons == nil {
		return nil, nil, nil, fmt.Errorf("Table key %s not found by manager but should be found", TPerson)
	}

	addresses := model.TableByKey(TAddress)

	if addresses == nil {
		return nil, nil, nil, fmt.Errorf("Table key %s not found by manager but should be found", TAddress)
	}

	return model, persons, addresses, nil

}
