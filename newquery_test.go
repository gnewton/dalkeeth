package dalkeeth

import (
	"testing"
)

func TestNewQuery_T1(t *testing.T) {
	setupTest()
	//mgr, err := initTestTables()
	_, err := initTestTables()
	if err != nil {
		t.Fatal(err)
	}
	// end setup
}
