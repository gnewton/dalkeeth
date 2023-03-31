package dalkeeth

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"strconv"
	"testing"
)

var ShouldHaveFailed = errors.New("Should have failed.")

func TestManager_Table_EmptyString(t *testing.T) {
	mgr, err := initTestTables()
	if err != nil {
		t.Fatal(err)
	}
	// end setup

	tName := ""
	hosts := mgr.Table(tName)
	if hosts != nil {
		t.Fatal(fmt.Errorf("Table with empty string key"))
	}

}

func TestManager_Table_UnknownString(t *testing.T) {
	mgr, err := initTestTables()
	if err != nil {
		t.Fatal(err)
	}
	// end setup

	tName := "not-existing-table-foo"
	hosts := mgr.Table(tName)
	if hosts != nil {
		t.Fatal(fmt.Errorf("Table should not exist: %s", tName))
	}

}

func TestManager_AddTable_KeyEmptyString(t *testing.T) {
	mgr, err := initTestTables()
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()

	newTable, err := NewTable("test1")
	if err != nil {
		t.Fatal(err)
	}
	// end setup
	err = mgr.AddTable("", newTable)
	if err == nil {
		t.Fatal(err)
	}
}

func TestManager_AddTable_NilTable(t *testing.T) {
	mgr, err := initTestTables()
	if err != nil {
		t.Fatal(err)
	}

	// end setup

	err = mgr.AddTable("valid", nil)
	if err == nil {
		t.Fatal(err)
	}
}

func TestManager_AddTable_KeyCollision(t *testing.T) {
	mgr, err := initTestTables()
	if err != nil {
		t.Fatal(err)
	}

	// end setup

	newTable, err := NewTable(TPerson)
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.AddTable(TPerson, newTable)
	if err == nil {
		t.Fatal(ShouldHaveFailed)
	}
}

func TestManager_CreateTablesSQL_NilDialect(t *testing.T) {
	mgr, err := initTestTables()
	if err != nil {
		t.Fatal(err)
	}

	// end setup

	mgr.dialect = nil

	_, err = mgr.CreateTablesSQL()
	if err == nil {
		t.Fatal(ShouldHaveFailed)
	}
}

func TestManager_CreateTableIndexesSQL_NilDialect(t *testing.T) {
	mgr, err := initTestTables()
	if err != nil {
		t.Fatal(err)
	}

	// end setup

	mgr.dialect = nil

	_, err = mgr.CreateTableIndexesSQL()
	if err == nil {
		t.Fatal(ShouldHaveFailed)
	}
}

func TestManager_AddForeignKey_NilTable(t *testing.T) {
	//t.Fatal(NotImplemented)
}

func Test_ManagerInitTables2(t *testing.T) {
	_, err := initTestTables()
	if err != nil {
		t.Fatal(err)
	}
}

func addForeignKey_Setup() (*Manager, *Table, *Table, error) {
	mgr, err := initTestTables()

	if err != nil {
		return nil, nil, nil, err
	}

	persons := mgr.Table(TPerson)

	if persons == nil {
		return nil, nil, nil, fmt.Errorf("Table key %s not found by manager but should be found", TPerson)
	}

	addresses := mgr.Table(TAddressK)

	if addresses == nil {
		return nil, nil, nil, fmt.Errorf("Table key %s not found by manager but should be found", TAddressK)
	}

	return mgr, persons, addresses, nil

}

func Test_Manager_AddForeignKey_UnknownForeignKeyField(t *testing.T) {
	mgr, persons, addresses, err := addForeignKey_Setup()

	if err != nil {
		t.Error(err)
	}

	if mgr.AddForeignKey(persons, "foo", addresses, FId) == nil {
		t.Fatal(fmt.Errorf("Failed identifying incorrect field"))
	}
}

func Test_Manager_AddForeignKey_UnknownForeignKeyFieldOtherField(t *testing.T) {
	mgr, err := initTestTables()

	if err != nil {
		t.Error(err)
	}

	persons := mgr.Table(TPerson)

	if persons == nil {
		t.Log(mgr.tablesMap)
		t.Fatal(fmt.Errorf("Table key %s not found by manager but should be found", TPerson))
	}

	addresses := mgr.Table(TAddressK)

	if addresses == nil {
		t.Fatal(fmt.Errorf("Table key %s not found by manager but should be found", TAddressK))
	}

	if mgr.AddForeignKey(persons, FId, addresses, "foo") == nil {
		t.Fatal(fmt.Errorf("Failed identifying incorrect field"))
	}
}

func Test_Manager_Manager_SaveTx(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	mgr, err := initAndWriteTestTables()
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()

	persons := mgr.Table(TPerson)
	if persons == nil {
		t.Error(errors.New("Persons cannot be found: is nil"))
	}

	mgr.dialect = new(DialectSqlite3)

	pk := int64(23)
	rec := persons.NewRecord()
	err = rec.AddValue(FId, pk)
	if err != nil {
		t.Fatal(err)
	}
	err = rec.AddValue(FName, "Fred")
	if err != nil {
		t.Fatal(err)
	}

	err = rec.AddValue(FAge, 54)
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.Begin()
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.SaveTx(rec)
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.Commit()
	if err != nil {
		t.Fatal(err)
	}

	valid, err := contains(mgr.db, rec.table.name, pk)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal(errors.New("Value not in database"))
	}

}

func Test_Manager_Manager_Save(t *testing.T) {
	mgr, err := initAndWriteTestTables()
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()

	persons := mgr.Table(TPerson)
	if persons == nil {
		t.Fatal(fmt.Errorf("Table key %s not found by manager but should be found", TPerson))
	}

	mgr.dialect = new(DialectSqlite3)

	pk := int64(23)
	rec := persons.NewRecord()
	err = rec.AddValue(FId, pk)
	if err != nil {
		t.Fatal(err)
	}
	err = rec.AddValue(FName, "Fred")
	if err != nil {
		t.Fatal(err)
	}

	err = rec.AddValue(FAge, 54)
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.Save(rec)
	if err != nil {
		t.Fatal(err)
	}

	// See if we can read the record that was just written
	valid, err := contains(mgr.db, rec.table.name, pk)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal(errors.New("Value not in database"))
	}
}

func Test_Manager_Manager_Batch(t *testing.T) {
	mgr, err := initAndWriteTestTables()
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()

	persons := mgr.Table(TPerson)
	if persons == nil {
		t.Fatal(fmt.Errorf("Table key %s not found by manager but should be found", TPerson))
	}

	mgr.dialect = new(DialectSqlite3)

	records, err := twoPersonRecords(persons)
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.Begin()
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.Batch(records)
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.Commit()
	if err != nil {
		t.Fatal(err)
	}

	// See if the 2 records added are readable
	valid, err := contains(mgr.db, records[0].table.name, VPersonID0)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal(errors.New("Value not in database"))
	}
	valid, err = contains(mgr.db, records[1].table.name, VPersonID1)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal(errors.New("Value not in database"))
	}
}

func Test_Manager_Manager_Begin_DBNil(t *testing.T) {
	mgr := NewManager()

	err := mgr.Begin()
	if err == nil {
		t.Fatal(ShouldHaveFailed)
	}
}

func Test_Manager_Manager_Begin_DoubleTx(t *testing.T) {
	mgr, err := initAndWriteTestTables()
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()

	err = mgr.Begin()
	if err != nil {
		t.Fatal(err)
	}

	// Should fail
	err = mgr.Begin()
	if err == nil {
		t.Fatal(ShouldHaveFailed)
	}
	t.Log(err)
}

func Test_Manager_Manager_Commit_NilTx(t *testing.T) {
	mgr, err := initAndWriteTestTables()
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()

	// Should fail
	err = mgr.Commit()
	if err == nil {
		t.Fatal(ShouldHaveFailed)
	}
	t.Log(err)
}

func Test_Manager_Manager_BatchMany(t *testing.T) {
	mgr, err := initAndWriteTestTables()
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()

	persons := mgr.Table(TPerson)
	if persons == nil {
		t.Fatal(fmt.Errorf("Table key %s not found by manager but should be found", TPerson))
	}

	records, err := nPersonRecords(persons, 10000)
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.Begin()
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.Batch(records)
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.Commit()
	if err != nil {
		t.Fatal(err)
	}

	valid, err := contains(mgr.db, records[0].table.name, VPersonID0)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal(errors.New("Value not in database"))
	}
	log.Println("Found", VPersonID0, "in DB")

	valid, err = contains(mgr.db, records[1].table.name, VPersonID1)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal(errors.New("Value not in database"))
	}
	log.Println("Found", VPersonID1, "in DB")
}

func contains(db *sql.DB, tableName string, id int64) (bool, error) {
	q := "SELECT id from " + tableName + " where id=?"
	var value int64

	row := db.QueryRow(q, id)
	err := row.Scan(&value)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return value == id, nil
}

func Test_Manager_Manager_Save_MissingNotNullValue(t *testing.T) {
	db, err := openTestDB()
	if err != nil {
		t.Error(err)
	}
	mgr, persons, _, err := addForeignKey_Setup()

	if err != nil {
		t.Error(err)
	}
	mgr.db = db
	mgr.dialect = new(DialectSqlite3)

	rec := persons.NewRecord()
	err = rec.AddValue(FId, 23)
	if err != nil {
		t.Fatal(err)
	}
	err = rec.AddValue(FAge, 54)
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.Begin()
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.Save(rec)
	if err == nil {
		t.Fatal(ShouldHaveFailed)
	}
}

func Test_Manager_Manager_Save_MissingPK(t *testing.T) {
	db, err := openTestDB()
	if err != nil {
		t.Error(err)
	}
	mgr, persons, _, err := addForeignKey_Setup()

	if err != nil {
		t.Error(err)
	}
	mgr.db = db
	mgr.dialect = new(DialectSqlite3)

	rec := persons.NewRecord()

	err = rec.AddValue(FName, "Fred")
	if err != nil {
		t.Fatal(err)
	}

	err = rec.AddValue(FAge, 54)
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.Begin()
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.Save(rec)
	if err == nil {
		t.Fatal(ShouldHaveFailed)
	}
}

func Test_Manager_Manager_Get(t *testing.T) {
	mgr, err := initAndWriteTestTables()
	if err != nil {
		t.Fatal(err)
	}
	defer mgr.Close()

	persons := mgr.Table(TPerson)
	if persons == nil {
		t.Fatal(fmt.Errorf("Table key %s not found by manager but should be found", TPerson))
	}

	mgr.dialect = new(DialectSqlite3)

	records, err := twoPersonRecords(persons)
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.Begin()
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.Batch(records)
	if err != nil {
		t.Fatal(err)
	}

	err = mgr.Commit()
	if err != nil {
		t.Fatal(err)
	}

	valid, err := contains(mgr.db, records[0].table.name, VPersonID0)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal(errors.New("Value not in database"))
	}
	log.Println("Found", VPersonID0, "in DB")

	valid, err = contains(mgr.db, records[1].table.name, VPersonID1)
	if err != nil {
		t.Fatal(err)
	}
	if !valid {
		t.Fatal(errors.New("Value not in database"))
	}
	log.Println("Found", VPersonID1, "in DB")

	// end setup

	rec, err := mgr.Get(persons, VPersonID0)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("VALUE=", rec.values[0].value)
	t.Log("VALUE=", rec.values[0].value)

	switch v := rec.values[0].value.(type) {
	case *int64:
		t.Log("is int", v)
	}
	//t.Log(rec.values[0].value)
	//t.Log(rec.values[1].value)
	//t.Log(rec.values[2].value)

}

// Helpers

func nPersonRecords(persons *Table, n int) ([]*Record, error) {
	records := make([]*Record, n)

	for i := 0; i < n; i++ {
		rec := persons.NewRecord()
		err := simplePersonRecord(rec, i)
		if err != nil {
			return nil, err
		}
		records[i] = rec
	}
	return records, nil
}

func simplePersonRecord(rec *Record, n int) error {
	rec.AddValue(FId, n)

	err := rec.AddValue(FName, "Fred_"+strconv.Itoa(n))
	if err != nil {
		return err
	}

	err = rec.AddValue(FAge, 54)
	if err != nil {
		return err
	}
	return nil
}

func twoPersonRecords(persons *Table) ([]*Record, error) {
	records := make([]*Record, 2)

	rec1 := persons.NewRecord()
	records[0] = rec1
	err := rec1.AddValue(FId, VPersonID0)
	if err != nil {
		return nil, err
	}
	err = rec1.AddValue(FName, "Fred")
	if err != nil {
		return nil, err
	}

	err = rec1.AddValue(FAge, 54)
	if err != nil {
		return nil, err
	}

	rec2 := persons.NewRecord()
	records[1] = rec2
	err = rec2.AddValue(FId, VPersonID1)
	if err != nil {
		return nil, err
	}
	err = rec2.AddValue(FName, "Harry")
	if err != nil {
		return nil, err
	}

	err = rec2.AddValue(FAge, 21)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func openTestDB() (*sql.DB, error) {
	return sql.Open("sqlite3", ":memory:")
}
