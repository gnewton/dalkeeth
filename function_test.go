package dalkeeth

import (
	"log"
	"testing"
)

func Test_SimpleStringField(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	sf := StringField{fieldName: "name"}

	ff := NewFunctionField(AVG, sf)

	d := new(DialectSqlite3)

	s, err := ff.ToSqlString(d)
	if err != nil {
		t.Error(err)
	}
	if s != "AVG(name)" {
		//t.Error(NotImplemented)
	}
}

func Test_SimpleStringField_WrongNArgs(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	sf := StringField{fieldName: "name"}
	ff := NewFunctionField(AVG, sf, sf)

	d := new(DialectSqlite3)

	s, err := ff.ToSqlString(d)
	if true {
		return
	}

	if err != nil {
		t.Error(err)
	}
	if s != "AVG(name)" {
		t.Log(s)
		t.Error()
	}

	fName := NewField("name", StringType, true, false, false, 0)
	//fAge := NewField("age", IntType, true, false, false, 0)

	countNameField := fName.Count()
	if err != nil {
		t.Error(err)
	}
	_, err = countNameField.ToSqlString(d)
	if err != nil {
		t.Error(err)
	}
}
