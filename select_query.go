package dalkeeth

type Ordering int

const (
	NoOrdering Ordering = iota
	ASC
	DESC
)

func (ord Ordering) String() string {
	return [...]string{"", "ASC", "DESC"}[ord]
}

type SelectQuery struct {
	Distinct       bool
	Fields         []AField
	rawFields      []string
	Pks            []int64
	From           []*Table
	FromRaw        []string
	Where          Condition
	GroupBy        []*Field
	Having         Condition
	Limit          int64
	Offset         int64
	OrderByFields  []*FieldOrdered
	GlobalOrdering Ordering
	SelectLimit    int64
	//
	validated bool
}

type FieldOrdered struct {
	Field
	function string
	as       string
	ordering Ordering
}

func (q *SelectQuery) Validated() bool {
	return q.validated
}

func (q *SelectQuery) Validate(m *Model) error {
	if err := zeroLength(q.Fields, "SelectQuery.fields"); err != nil {
		return err
	}

	err := rawFieldsToFields(q, m)
	if err != nil {
		return err
	}

	err = rawFromTablesToTables(q, m)
	if err != nil {
		return err
	}

	q.validated = true
	return nil
}

// ////////////////////////////////////
// Run
func (q *SelectQuery) First() (*InRecord, error) {
	return nil, nil
}
func (q *SelectQuery) Last() (*InRecord, error) {
	return nil, nil
}

func (q *SelectQuery) Rows() (*Rows, error) {
	return nil, nil
}

func (q *SelectQuery) Exists() (bool, error) {
	return false, nil
}

func (q *SelectQuery) Pluck() ([]any, error) {
	return nil, nil
}

type Rows struct { //temporary
}

//////////////////// ex2 NEW

func (q *SelectQuery) Select(fields ...AField) {
	for i := 0; i < len(fields); i++ {
		//q.fields = append(q.fields, fields[i])
	}
}
