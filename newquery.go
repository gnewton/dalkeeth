package dalkeeth

type Join2 struct {
	f1, f2 *Field
}

type OrderBy struct {
	field    AField
	ordering int // ASC, DESC
}

type Query struct {
	selectFields []AField
	selectAny    []string
	whereEquals  []AField
	where        []*Condition
	joins        []*Join2
	joinsByName  map[string]*Join2
	groupBy      []*Field // Can this be AField?
	having       *Condition
	orderBy      []*OrderBy
	offset       int64
	limit        int64
}

func NewQuery() *Query {
	return &Query{
		selectFields: make([]AField, 0),
		selectAny:    make([]string, 0),
		whereEquals:  make([]AField, 0),
		where:        make([]*Condition, 0),
		joins:        make([]*Join2, 0),
		groupBy:      make([]*Field, 0),
		orderBy:      make([]*OrderBy, 0),
		offset:       -1,
		limit:        -1,
	}
}

func (q *Query) Select(fields ...AField) *Query {
	for i := 0; i < len(fields); i++ {
		q.selectFields = append(q.selectFields, fields[i])
	}
	return q
}

func (q *Query) SelectAny(strs ...string) *Query {
	for i := 0; i < len(strs); i++ {
		q.selectAny = append(q.selectAny, strs[i])
	}
	return q
}

func (q *Query) Join(jf1, jf2 *Field) *Query {
	q.joins = append(q.joins, &Join2{f1: jf1, f2: jf2})
	return nil
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
