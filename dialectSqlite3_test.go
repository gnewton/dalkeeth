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
	mgr, err := initAndPopulateTestTables()

	if err != nil {
		t.Error(err)
	}

	_, err = mgr.dialect.ExtractTable(mgr.db, TPerson)
	if err != nil {
		t.Error(err)
	}

}
