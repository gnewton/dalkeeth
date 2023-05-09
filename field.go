package dalkeeth

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
)

type AField interface {
	ToSqlString(Dialect) (string, error)
}

type StringField struct {
	fieldName string
}

func (sf StringField) ToSqlString(d Dialect) (string, error) {
	return sf.fieldName, nil
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

type FieldAs struct {
	field *Field
	alias string
}

func (fa *FieldAs) ToSqlString(d Dialect) (string, error) {
	return d.FieldAsSql(fa)
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

func (f *Field) ToSqlString(d Dialect) (string, error) {
	return f.name, nil
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

func (f *Field) As(alias string) AField {
	return new(FieldAs)
}

// SQL functions used in queries

func (f *Field) Count() AField {
	return NewFunctionField(COUNT, f)
}

func (f *Field) Avg() AField {
	return nil
}

func (f *Field) Round() AField {
	return nil
}

// Conditions on fields (where, having)

func (f *Field) In(in ...any) *Condition {
	return nil
}

func (f *Field) Is(v any) *Condition {
	return nil
}

func (f *Field) IsGreaterThan(v any) *Condition {
	return nil
}

func (f *Field) IsLessThan(v any) Condition {
	return nil
}
