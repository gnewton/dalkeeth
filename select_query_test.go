package dalkeeth

import (
	"testing"
)

func Test00(t *testing.T) {
	age := Field{name: "age", fieldType: IntType}
	//ageSelect := &SelectField{Field: age, function: "MAX", as: "MaxValue"}
	ageSelect := age.SelectFieldFuncAs("MAX", "MaxValue")
	name := Field{name: "name", fieldType: StringType}
	nameSelect := name.SelectField()

	q := SelectQuery{
		distinct: true,
		fields:   []*SelectField{ageSelect, nameSelect},
		pks:      []int64{54, 767},
		where:    WN(nameSelect, IsNotNull),
		groupBy:  []*Field{&name},
		having:   W(&age, GT, 100),
		offset:   1200,
		limit:    100,
		ordering: ASC,
	}

	if err := q.Validate(); err != nil {
		t.Error(err)
	}

}
