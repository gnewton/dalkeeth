package dalkeeth

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)

// support
// w := NewWhere(S("a=b")).And(SP("foo=", "?")).Or(FP(field2, Equals, "?")).Or(Not(SC(field1, Equals, field2)))w

type Operator int
type LogicalOperator int

const (
	LOr LogicalOperator = iota
	LAnd
	LNot
)

const (
	Between Operator = iota
	EQ
	GE
	GT
	In
	IsNotNull
	IsNotTrue
	IsNull
	IsTrue
	LE
	LT
	Like
	NE
	NotBetween
	NotIn
)

var NumberOperators = []Operator{
	Between, EQ, GE, GT, LT, LE, NE, NotBetween,
}

var StringOnlyOperators = []Operator{
	Like,
}

func (op LogicalOperator) String() string {
	switch op {
	case LOr:
		return " OR "
	case LAnd:
		return " AND "

	case LNot:
		return " NOT "

	}
	return "unknown"
}

func (op Operator) String() string {
	switch op {
	case Between:
		return " BETWEEN "
	case EQ:
		return " = "

	case GE:
		return " >= "

	case GT:
		return " > "

	case In:
		return " IN "

	case IsNotNull:
		return " IS NOT NULL "

	case IsNotTrue:
		return " IS NOT TRUE "

	case IsNull:
		return " IS NULL "

	case IsTrue:
		return " IS TRUE "

	case LE:
		return " < "

	case LT:
		return " > "

	case Like:
		return " LIKE "

	case NE:
		return " <> "

	case NotBetween:
		return " NOT BETWEEN "

	case NotIn:
		return " NOT IN "
	}

	return "unknown"
}

const MaxInLength = 100

func (op Operator) ArgStart() string {
	switch op {
	case In, NotIn:
		return "["
	}
	return ""
}

func (op Operator) ArgEnd() string {
	switch op {
	case In, NotIn:
		return "]"
	}
	return ""
}

func (op Operator) MinMaxArgs() (int, int) {
	switch op {
	case Between, NotBetween:
		return 2, 2
	case In, NotIn:
		return 1, MaxInLength
	case IsNotNull, IsNotTrue, IsNull, IsTrue:
		return 0, 0
	}
	return 1, 1
}

func (op Operator) StringValueOperator() bool {
	switch op {
	case Between, NotBetween, IsNotNull, IsNotTrue, IsNull, IsTrue:
		return false
	}
	return true
}

type Condition interface {
	Evaluate(depth int) (string, error)
	Validate() error
	String() string // Not sql; just for debugging
}

type LHS interface {
	string | *Field
}

type Numbers interface {
	int | int64 | float64
}

type Values interface {
	//int | int64 | float64 | string | *Field
	Numbers | string | *Field | *StringField
}

var None = []int{}

func WN[L LHS](left L, op Operator) Condition {
	return W(left, op, None...)
}

func W[L LHS, V Values](left L, op Operator, values ...V) Condition {
	switch l := any(left).(type) {
	case string:
		switch v := any(values).(type) {
		// "Native"
		case []string:
			e := new(SP_GenericCondition[string, string])
			e.left = l
			e.values = v
			e.op = op
			e.valuesAreStrings = true
			return e
		case []float64:
			e := new(SP_GenericCondition[string, float64])
			e.left = l
			e.values = v
			e.op = op
			return e
		case []*Field:
			e := new(SP_GenericCondition[string, *Field])
			e.left = l
			e.values = v
			e.op = op
			return e
		case []int64:
			e := new(SP_GenericCondition[string, int64])
			e.left = l
			e.values = v
			e.op = op
			return e
		// Types mapped to "Native"
		case []int:
			e := new(SP_GenericCondition[string, int64])
			e.left = l
			e.values = toInt64(v)
			e.op = op
			return e
		}
	case *Field:
		switch v := any(values).(type) {
		case []int:
			e := new(SP_GenericCondition[*Field, int64])
			e.left = l
			e.values = toInt64(v)
			e.op = op
			return e
		case []string:
			e := new(SP_GenericCondition[*Field, string])
			e.left = l
			e.values = v
			e.op = op
			e.valuesAreStrings = true
			return e
		case []float64:
			e := new(SP_GenericCondition[*Field, float64])
			e.left = l
			e.values = v
			e.op = op
			return e
		case []*Field:
			e := new(SP_GenericCondition[*Field, *Field])
			e.left = l
			e.values = v
			e.op = op
			return e
		case []int64:
			e := new(SP_GenericCondition[*Field, int64])
			e.left = l
			e.op = op
			e.values = v
			return e
		}
	}
	return nil
}

type SP_GenericCondition[L LHS, V Values] struct {
	left             L
	op               Operator
	values           []V
	valuesAreStrings bool
	validated        bool
}

func toInt64(v []int) []int64 {
	i64 := make([]int64, len(v))
	for i := 0; i < len(v); i++ {
		i64[i] = int64(v[i])
	}
	return i64
}

func (expr *SP_GenericCondition[Left, Values]) Validate() error {
	expr.validated = true
	log.Println("SP_GenericCondition[Left, Values]) Validate---->", expr.validated)
	return nil
}

func (expr *SP_GenericCondition[Left, Values]) String() string {
	return "SP_GenericCondition: " + expr.op.String()
}

func (expr *SP_GenericCondition[Left, Values]) Evaluate(d int) (string, error) {
	log.Println("SP_GenericCondition[Left, Values]) Evaluate---->", expr.validated)
	if !expr.validated {
		return "", errors.New("SP_GenericCondition.Evaluate: Cannot Evaluate without first being Validated:" + expr.String())
	}

	var e string

	err := validateOperationWithValuesCount(expr.op, len(expr.values))
	log.Println(expr.left, expr.values)

	if err != nil {
		return "", err
	}

	// if OP is LIKE, expr.values MUST be string
	if !expr.valuesAreStrings && expr.op == Like {
		return "", fmt.Errorf("Like must be with string field %v", expr)
	}

	if expr.valuesAreStrings && !expr.op.StringValueOperator() {
		return "", fmt.Errorf("String values need string operator; have %s", expr.op.String())
	}

	switch l := any(expr.left).(type) {
	case string:
		e += l + expr.op.String() + expr.op.ArgStart()

		rawValues, err := toValues(expr.values)
		if err != nil {
			return "", err
		}
		e += rawValues
		e += expr.op.ArgEnd()
		return e, nil
	case *Field:
		e += l.String() + expr.op.String() + expr.op.ArgStart()
		rawValues, err := toValues(expr.values)
		if err != nil {
			return "", err
		}
		e += rawValues
		e += expr.op.ArgEnd()
		return e, nil
	default:
		return "Field not implemented", nil
	}
}

func toValues[V Values](values []V) (string, error) {
	var s string
	for i := 0; i < len(values); i++ {
		if i > 0 {
			s += ", "
		}
		v := values[i]
		raw, err := valueToString(v)
		if err != nil {
			return "", err
		}
		s += raw
	}
	return s, nil
}

func valueToString[V Values](value V) (string, error) {
	switch v := any(value).(type) {
	case int64:
		return strconv.FormatInt(v, 10), nil
	case float64:
		return fmt.Sprintf("%f", v), nil
	case string:
		return "\"" + v + "\"", nil
	case *Field:
		return v.String(), nil
	}
	return "", errors.New("Unknown type")
}

func (f *Field) String() string {
	return "<field>"
}

func Not(e Condition) Condition {
	return &NotCondition{
		e: e,
	}

}

func Evaluate(e Condition) (string, error) {
	return e.Evaluate(0)
}

type NotCondition struct {
	e         Condition
	validated bool
}

func (expr NotCondition) Validate() error {
	expr.validated = true
	return nil
}

func (expr NotCondition) String() string {
	return "NotCondition"
}

func (n *NotCondition) Evaluate(d int) (string, error) {
	eval, err := n.e.Evaluate(d + 1)
	if err != nil {
		return "", err
	}
	return LNot.String() + eval, nil
}

type LogicalCondition struct {
	exp1, exp2 Condition
	exps       []Condition
	op         LogicalOperator
	validated  bool
}

func (expr LogicalCondition) Validate() error {
	expr.validated = true
	return nil
}

func (expr LogicalCondition) String() string {
	return "LogicalCondition"
}

func (o LogicalCondition) Evaluate(d int) (string, error) {
	var s string

	if o.exp1 == nil {
		return "", errors.New("Left expression is nil")
	}

	if o.exp2 == nil {
		return "", errors.New("Right expression is nil")
	}

	err := o.exp1.Validate()
	if err != nil {
		return "", err
	}
	err = o.exp2.Validate()
	if err != nil {
		return "", err
	}

	eval, err := o.exp1.Evaluate(d + 1)
	if err != nil {
		return "", err
	}
	s += eval + o.op.String()

	eval, err = o.exp2.Evaluate(d + 1)
	if err != nil {
		return "", err
	}
	s += eval

	for i := 0; i < len(o.exps); i++ {
		err = o.exps[i].Validate()
		if err != nil {
			return "", err
		}
		eval, err = o.exps[i].Evaluate(d + 1)
		if err != nil {
			return "", err
		}
		s += o.op.String() + eval
	}
	if d > 0 {
		s = "(" + s + ")"
	}
	return s, nil
}

func Or(exp1, exp2 Condition, exps ...Condition) Condition {
	return LogicalCondition{
		exp1: exp1,
		exp2: exp2,
		exps: exps,
		op:   LOr,
	}
}

func And(exp1, exp2 Condition, exps ...Condition) Condition {
	return LogicalCondition{
		exp1: exp1,
		exp2: exp2,
		exps: exps,
		op:   LAnd,
	}
}

func validateOperationWithValuesCount(op Operator, l int) error {
	min, max := op.MinMaxArgs()
	log.Println(op.String(), min, max)
	if l < min {
		return fmt.Errorf("Too few arguments %d for operator %s: need > %d", l, op.String(), min)
	}

	if l > max {
		return fmt.Errorf("Too many arguments %d for operator %s: need < %d", l, op.String(), max)
	}

	return nil

}

//////////////
