package dalkeeth

type Join2 struct {
	f1, f2 *Field
}

type OrderBy struct {
	field    AField
	ordering int // ASC, DESC
}

type Query struct {
	selectFields []*Field
	selectRaw    []string
	fromTables   []*Table
	fromRaw      string
	whereEquals  []AField
	where        []*Condition
	whereRaw     string
	joins        []*Join2
	joinsByName  map[string]*Join2
	groupBy      []*Field // Can this be AField?
	having       *Condition
	orderBy      []*OrderBy
	offset       int64
	limit        int64
	initialized  bool
}

func NewQuery() *Query {
	return queryInitialize()
}

func queryInitialize() *Query {
	return &Query{
		selectFields: make([]*Field, 0),
		selectRaw:    make([]string, 0),
		whereEquals:  make([]AField, 0),
		where:        make([]*Condition, 0),
		joins:        make([]*Join2, 0),
		groupBy:      make([]*Field, 0),
		orderBy:      make([]*OrderBy, 0),
		offset:       -1,
		limit:        -1,
		initialized:  true,
	}
}

func (q *Query) From(tables ...*Table) *Query {
	for i := 0; i < len(tables); i++ {
		q.fromTables = append(q.fromTables, tables[i])
	}
	return q
}

func (q *Query) FromRaw(s string) *Query {
	q.fromRaw = s
	return q
}

func (q *Query) Select(fields ...*Field) *Query {
	for i := 0; i < len(fields); i++ {
		q.selectFields = append(q.selectFields, fields[i])
	}
	return q
}

func (q *Query) SelectByName(strs ...string) *Query {
	for i := 0; i < len(strs); i++ {
		q.selectRaw = append(q.selectRaw, strs[i])
	}
	return q
}

func (q *Query) Join(jf1, jf2 *Field) *Query {
	q.joins = append(q.joins, &Join2{f1: jf1, f2: jf2})
	return nil
}

func (q *Query) WhereRaw(s string) *Query {
	q.whereRaw = s
	return q
}

func (q *Query) Where(cs ...*Condition) *Query {
	for i := 0; i < len(cs); i++ {
		q.where = append(q.where, cs[i])
	}
	return q
}

func (q *Query) WhereEquals(jf1, jf2 *Field) *Query {
	q.whereEquals = append(q.whereEquals, jf1, jf2)
	return q
}

func (q *Query) GroupBy(fields ...*Field) *Query {
	for i := 0; i < len(fields); i++ {
		q.groupBy = append(q.groupBy, fields[i])
	}
	return q
}

func (q *Query) Having(c *Condition) *Query {
	q.having = c
	return q
}

func (q *Query) OrderBy(ob ...*OrderBy) *Query {
	for i := 0; i < len(ob); i++ {
		q.orderBy = append(q.orderBy, ob[i])
	}
	return q
}

func (q *Query) Offset(offset int64) *Query {
	q.offset = offset
	return q
}

func (q *Query) Limit(limit int64) *Query {
	q.limit = limit
	return q
}

func (q *Query) Freeze() (*Query, error) {
	return nil, NotImplemented
}
