package dalkeeth

import (
	//"database/sql"
	//"errors"
	//"fmt"
	"log"
	"testing"
)

func TestDialectSqlite3(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	setupTest()
	mdl0, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}
	sess, err := writeTestModelSchema(mdl0)

	if err != nil {
		t.Fatal(err)
	}

	_, err = sess.dialect.ExtractTable(sess.db, TPerson)
	if err != nil {
		t.Error(err)
	}

}
