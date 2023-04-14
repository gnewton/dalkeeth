package dalkeeth

import (
	//"database/sql"
	//"errors"
	//"fmt"
	"testing"
)

func TestDialectSqlte3(t *testing.T) {
	setupTest()
	mgr, err := initAndWriteTestTables()

	if err != nil {
		t.Error(err)
	}

	_, err = mgr.dialect.ExtractTable(mgr.db, TPerson)
	if err != nil {
		t.Error(err)
	}

}
