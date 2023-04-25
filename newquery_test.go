package dalkeeth

import (
	"testing"
)

func TestNewQuery_T1(t *testing.T) {
	setupTest()
	//mgr, err := defineTestModel()
	_, err := defineTestModel()
	if err != nil {
		t.Fatal(err)
	}
	// end setup
}
