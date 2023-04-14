package dalkeeth

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if true {
		log.SetOutput(ioutil.Discard)
	}
	exitVal := m.Run()
	os.Exit(exitVal)
}

func setupTest() {

}
