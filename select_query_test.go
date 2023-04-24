package dalkeeth

import (
	"testing"
)

func Test00(t *testing.T) {
	age := Field{name: "age", fieldType: IntType}
	name := Field{name: "name", fieldType: StringType}

	q := SelectQuery{
		fields:         []*Field{&age, &name},
		pks:            []int64{54, 767},
		where:          WN(&name, IsNotNull),
		groupBy:        []*Field{&name},
		having:         W(&age, GT, 100),
		offset:         1200,
		limit:          100,
		globalOrdering: ASC,
		orderByFields:  []*FieldOrdered{},
	}

	if err := q.Validate(); err != nil {
		t.Error(err)
	}

}
