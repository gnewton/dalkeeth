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
	sess, err := initAndWriteTestTableSchema()

	if err != nil {
		t.Error(err)
	}

	_, err = sess.dialect.ExtractTable(sess.db, TPerson)
	if err != nil {
		t.Error(err)
	}

}
