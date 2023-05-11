package dalkeeth

import (
	"testing"
)

func TestNewQuery_T1(t *testing.T) {
	setupTest()
	//mgr, err := testModel0()
	_, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}
	// end setup
}
