package dalkeeth

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
)

func TestSqlCreate(t *testing.T) {
	tbl, err := testTable()

	if err != nil {
		t.Error(err)
	}

	_, err = tbl.CreateTableSql()
	if err != nil {
		t.Error(err)
	}
}

func TestIndexSql(t *testing.T) {
	tbl, err := testTable()
	if err != nil {
		t.Error(err)
	}

	index := Index{
		table: tbl,
		fields: []*Field{
			tbl.fields[1],
		},
	}
	_, err = index.CreateSql(0)
	if err != nil {
		t.Error(err)
	}
}

func TestIndexSqlWithDB(t *testing.T) {
	db, tbl, err := simpleTestTable()
	if err != nil {
		t.Error(err)
	}

	// end setup

	index := Index{
		table: tbl,
		fields: []*Field{
			tbl.fields[1],
		},
	}
	createIndexSql, err := index.CreateSql(0)

	if err != nil {
		t.Error(err)
	}

	_, err = db.Exec(createIndexSql)
	if err != nil {
		t.Log(createIndexSql)
		t.Error(err)
	}
}

func TestSqlWithDB(t *testing.T) {
	tbl, err := testTable()
	if err != nil {
		t.Error(err)
	}
	createSql, err := tbl.CreateTableSql()

	if err != nil {
		t.Error(err)
	}

	db, err := openTestDB()
	if err != nil {
		t.Error(err)
	}

	_, err = db.Exec(createSql)
	if err != nil {
		t.Log(createSql)
		t.Error(err)
	}

}

func TestSql_NoTableName(t *testing.T) {
	tbl, err := testTable()
	if err != nil {
		t.Error(err)
	}
	tbl.name = ""
	_, err = tbl.CreateTableSql()

	if err == nil {
		t.Error(err)
	}
}

func TestSql_NoFields(t *testing.T) {
	tbl, err := testTable()
	if err != nil {
		t.Error(err)
	}
	tbl.fields = nil
	_, err = tbl.CreateTableSql()

	if err == nil {
		t.Error(err)
	}
}

func TestSql_MultiplePrimaryKeys(t *testing.T) {
	tbl, err := testTable()
	if err != nil {
		t.Error(err)
	}
	f := Field{
		name:      "test",
		fieldType: IntType,
		pk:        true,
	}
	_, err = tbl.AddField(&f)
	if err == nil {
		t.Error("This should be an error")
	}
}

func TestFindFieldValueById(t *testing.T) {
	db, tbl, err := simpleTestTable()
	if err != nil {
		t.Error(err)
	}

	// end setup
	var name string
	ok, err := FindFieldValueById(tbl, db, TestId0, tbl.fields[1], &name)
	if err != nil {
		t.Error(err)
	}

	if !ok {
		t.Fatal(fmt.Errorf("Failed finding by id=%d", TestId0))
	}

	if name != TestName0 {
		t.Fatal(fmt.Errorf("Failed finding by id=%d  Wrong string=%s", TestId0, name))
	}
}

const TestTable0 = "test_table0"

const TestId0 = 0
const TestId1 = 42

const TestName0 = "Fred"
const TestName1 = "Xavier"

func TestSaveRecord(t *testing.T) {
	_, _, err := simpleTestTable()
	if err != nil {
		t.Error(err)
	}
}

func simpleTestTable() (*sql.DB, *Table, error) {
	db, err := openTestDB()
	if err != nil {
		return nil, nil, err
	}
	tbl, err := testTable()
	if err != nil {
		return nil, nil, err
	}
	createSql, err := tbl.CreateTableSql()
	if err != nil {
		return nil, nil, err
	}
	_, err = db.Exec(createSql)
	if err != nil {
		log.Println(createSql)
		return nil, nil, err
	}

	err = populateTable(db, tbl)
	if err != nil {
		return nil, nil, err
	}
	return db, tbl, nil
}

func testTable() (*Table, error) {
	model := NewModel()
	tbl, err := model.NewTable(TestTable0)
	if err != nil {
		return nil, err
	}

	tbl.AddField(&Field{
		name:      "id",
		fieldType: IntType,
		pk:        true,
	})

	tbl.AddField(&Field{
		name:      "name",
		fieldType: StringType,
	})
	return tbl, nil
}

func populateTable(db *sql.DB, tbl *Table) error {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	rec := tbl.NewRecord()

	err := rec.SetValue("id", TestId0)
	if err != nil {
		log.Println(err)
		return err
	}
	err = rec.SetValue("name", TestName0)
	if err != nil {
		log.Println(err)
		return err
	}

	rec2 := tbl.NewRecord()

	err = rec2.SetValue("id", TestId1)
	if err != nil {
		log.Println(err)
		return err
	}
	err = rec2.SetValue("name", TestName1)
	if err != nil {
		log.Println(err)
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		log.Println(err)
		return err
	}
	err = rec.Save(tx)
	if err != nil {
		log.Println(err)
		return err
	}
	err = rec2.Save(tx)
	if err != nil {
		log.Println(err)
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return err
	}

	q := "SELECT id, name from " + tbl.name

	row := db.QueryRow(q)
	var id int
	var name string
	err = row.Scan(&id, &name)
	if err != nil {
		log.Println(q)
		return err
	}
	return nil

}
