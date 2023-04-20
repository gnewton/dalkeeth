package dalkeeth

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
)

type AField interface {
	ToSqlString(Dialect) string
}

type StringField struct {
	field string
}

func (sf *StringField) ToSqlString(d Dialect) string {
	return sf.field
}

type FieldType int

const (
	IntType FieldType = iota // int64
	StringType
	BoolType
	FloatType     // float64
	ByteArrayType // []byte
	//
	FunctionType
)

func (ft FieldType) String() string {
	return [...]string{"IntType", "StringType", "BoolType", "FloatType", "ByteArrayType", "FunctionType"}[ft]
}

type Field struct {
	name         string
	fieldType    FieldType
	pk           bool
	indexed      bool
	length       int
	unique       bool
	notNull      bool
	defaultValue string
	table        *Table
}

func validTypeForField(v any, f *Field) error {

	switch t := v.(type) {
	case int:
		if f.fieldType != IntType {
			return fmt.Errorf("Table %s Field %s: Value %d is int; field type is %s", f.table.name, f.name, t, f.fieldType)
		}

	case string:
		if f.fieldType != StringType {
			return fmt.Errorf("Table %s Field %s: Value %s is string; field type is %s", f.table.name, f.name, t, f.fieldType)
		}

	case sql.NullString:
		if f.fieldType != StringType {
			return fmt.Errorf("Table %s Field %s: Value %s is string; field type is %s", f.table.name, f.name, t.String, f.fieldType)
		}
	}

	return nil
}
func NewField(name string, fieldType FieldType, pk, indexed, notNull bool, length int) *Field {
	f := new(Field)
	f.name = name
	f.fieldType = fieldType
	f.pk = pk
	f.indexed = indexed
	f.notNull = notNull
	f.length = length
	return f
}

func (f *Field) SelectField() *SelectField {
	sf := &SelectField{Field: *f}
	return sf
}

func (f *Field) ToSqlString(d Dialect) string {
	log.Fatal(NotImplemented)
	return "Unimplemented"
}

func (f *Field) SelectFieldFuncAs(function, as string) *SelectField {
	sf := f.SelectField()
	sf.function = function
	sf.as = as

	return sf
}

func (f *Field) CreateFieldSql() (string, error) {
	if f.name == "" {
		return "", errors.New("Field name is empty")
	}

	s := f.name + " " + sqlFieldType(f)

	return s, nil
}

func sqlFieldType(f *Field) string {
	var s string

	switch f.fieldType {
	case IntType:
		s = "INT"
	case StringType:
		if f.length == 0 {
			s = "TEXT"
		} else {
			s = "varchar(" + strconv.Itoa(f.length) + ")"
		}
	case FloatType:
		s = "REAL"
	}
	if f.notNull {
		s += " NOT NULL"
	}
	if f.unique {
		s += " UNIQUE"
	}
	if f.pk {
		s += " PRIMARY KEY"
	}
	return s

}

// Used in queries; functions

func (f *Field) Count() AField {
	countField := NewFunctionField(COUNT, f)
	return countField
}

func (f *Field) Avg() AField {
	return nil
}

func (f *Field) Round() AField {
	return nil
}

// Conditions on fields (where, having)
func (f *Field) Is() *Condition {
	return nil
}

func (f *Field) In(in ...any) *Condition {
	return nil
}

func (f *Field) GreaterThan(v any) *Condition {
	return nil
}
