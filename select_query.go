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

type Fields []*Field

type SelectQuery struct {
	distinct    bool
	fields      []*SelectField
	pks         []int64
	where       Condition
	groupBy     []*Field
	having      Condition
	limit       int64
	offset      int64
	orderBy     []*SelectField
	ordering    Ordering
	validated   bool
	selectLimit int64
}

type SelectField struct {
	Field
	function string
	as       string
	ordering Ordering
}

func (q *SelectQuery) Validate() error {
	if err := zeroLength(q.fields, "SelectQuery.fields"); err != nil {
		return err
	}

	if err := zeroLength(q.fields, "SelectQuery.fields"); err != nil {
		return err
	}
	return nil
}

// Run
func (q *SelectQuery) First() (*Record, error) {
	return nil, nil
}
func (q *SelectQuery) Last() (*Record, error) {
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
