package dalkeeth

import (
	// "errors"
	"fmt"
	"testing"
)

func TestSimple_Positive(t *testing.T) {
	exp := W("MAX(foo)", NotBetween, 43, 55)

	err := exp.Validate()
	if err != nil {
		t.Error(err)
	}
	s, err := Evaluate(exp)

	if err != nil {
		t.Error(err)
	}
	fmt.Println(s)
}

func TestComplex_Positive(t *testing.T) {
	field1 := new(Field)
	field2 := new(Field)
	exp := Or(And(W("MAX(foo)", NotBetween, 43, 55), WN(field1, IsNotNull)), And(W(field1, GT, field2), W("m", LT, 54.5), W("name", Like, "smith"), W("name", In, "smith", "rogers")))

	err := exp.Validate()
	if err != nil {
		t.Error(err)
	}
	s, err := Evaluate(exp)

	if err != nil {
		t.Error(err)
	}
	fmt.Println(s)
}

func TestSimpleStringString_Positive(t *testing.T) {
	for i := Between; i <= NotIn; i++ {
		min, max := i.MinMaxArgs()
		// Take only 1 arg
		if min == 1 && max == 1 && i.StringValueOperator() {
			exp := W("name", i, "Fred%")

			err := exp.Validate()
			if err != nil {
				t.Error(err)
			}

			s, err := Evaluate(exp)

			if err != nil {
				t.Error(err)
			}
			fmt.Println(s)
		}
	}

}
