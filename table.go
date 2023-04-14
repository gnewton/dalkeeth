package dalkeeth

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
)

const ID = "id"

type Record struct {
	table     *Table
	values    []*Value
	valuesMap map[string]*Value
}

type Value struct {
	field    *Field
	value    any
	isSet    bool
	isWanted bool
}

type Table struct {
	name        string
	pk          *Field
	fields      []*Field
	fieldsMap   map[string]*Field
	indexes     []*Index
	foreignKeys []*ForeignKey
}

type ForeignKey struct {
	tbl          *Table
	field        *Field
	foreignTable *Table
	foreignKey   *Field
}

type Join struct {
	segments []JoinSegment
}

type JoinSegment struct {
	t1 *Table
	f1 *Field
	t2 *Table
	f2 *Field
}

// name: idx_table_f0_f1_...
type Index struct {
	fields []*Field
	unique bool
	table  *Table
}

func NewTable(name string) (*Table, error) {
	if name == "" {
		return nil, errors.New("Table name is empty")
	}
	tbl := new(Table)
	tbl.name = name
	tbl.fieldsMap = make(map[string]*Field, 0)
	return tbl, nil
}

func (rec *Record) GetInt(fieldName string, vv *int64) error {
	v, err := rec.Value(fieldName)
	if err != nil {
		return err
	}

	switch p := v.value.(type) {
	case *int64:
		vv = p
	default:
		return errors.New("Wrong type: expecting *int64")
	}

	return nil
}

func (rec *Record) GetString(fieldName string, vv *string) error {
	v, err := rec.Value(fieldName)
	if err != nil {
		return err
	}

	switch p := v.value.(type) {
	case *string:
		vv = p
		return nil
	default:
		return errors.New("Wrong type: expecting *string")
	}

}

func (rec *Record) GetBool(fieldName string, vv *bool) error {
	v, err := rec.Value(fieldName)
	if err != nil {
		return err
	}

	switch p := v.value.(type) {
	case *bool:
		vv = p
		return nil
	default:
		return errors.New("Wrong type: expecting *bool")
	}

}

func (rec *Record) GetFloat(fieldName string, vv *float64) error {
	v, err := rec.Value(fieldName)
	if err != nil {
		return err
	}

	switch p := v.value.(type) {
	case *float64:
		vv = p
		return nil
	default:
		return errors.New("Wrong type: expecting *float64")
	}

}

func (rec *Record) Value(fieldName string) (*Value, error) {
	if fieldName == "" {
		return nil, errors.New("Field name is empty")
	}
	var v *Value
	var ok bool

	if v, ok = rec.valuesMap[fieldName]; !ok {
		return nil, fmt.Errorf("Field: %s is not in table: %s", fieldName, rec.table.name)
	}
	if v.field.fieldType != StringType {
		return nil, fmt.Errorf("Field type is not int: field type is: %s in table: %s", v.field.fieldType, rec.table.name)
	}
	return v, nil
}

func (rec *Record) Save(tx *sql.Tx) error {
	insert := "INSERT INTO " + rec.table.name + "("
	for i := 0; i < len(rec.values); i++ {
		if i != 0 {
			insert += COMMA_SPACE
		}
		insert += rec.table.fields[i].name
	}
	insert += ") VALUES ("
	var valueArray []any

	for i := 0; i < len(rec.values); i++ {
		if i != 0 {
			insert += COMMA_SPACE
		}
		insert += "$" + strconv.Itoa(i+1)
		valueArray = append(valueArray, rec.values[i].value)
	}
	insert += ")"

	results, err := tx.Exec(insert, valueArray...)

	if err != nil {
		log.Println(insert, valueArray)
		return err
	}

	nrows, err := results.RowsAffected()
	if err != nil {
		log.Println(insert, valueArray)
		return err
	}

	if nrows != 1 {
		return fmt.Errorf("Expected only 1 row effected; %d rows affected", nrows)
	}
	return nil
}

func (rec *Record) AddValue(name string, value any) error {
	var v *Value
	var ok bool
	if v, ok = rec.valuesMap[name]; !ok {
		return fmt.Errorf("Field name not found:[%s] in table[%s]", name, rec.table.name)
	}

	err := validTypeForField(value, v.field)
	if err != nil {
		return err
	}

	v.value = &value
	v.isSet = true
	return nil
}

func (t *Table) NewRecord() *Record {
	rec := new(Record)
	rec.table = t
	if t == nil {
		log.Fatal("t == nil")
	}
	if t.fields == nil {
		log.Fatal("t.fields == nil")
	}
	rec.valuesMap = make(map[string]*Value, len(t.fields))
	rec.values = make([]*Value, len(t.fields))
	for i := 0; i < len(rec.values); i++ {
		val := new(Value)
		val.field = t.fields[i]
		rec.values[i] = val
		val.isWanted = true
		rec.valuesMap[t.fields[i].name] = val
	}
	return rec
}

func (t *Table) AddFields(fs ...*Field) error {
	if len(fs) == 0 {
		return errors.New("Empty array of *Fields")
	}
	for i := 0; i < len(fs); i++ {
		if _, err := t.AddField(fs[i]); err != nil {
			return err
		}
	}
	return nil

}

func (t *Table) AddField(f *Field) (*Field, error) {
	if f == nil {
		return nil, errors.New("AddField: Nil pointer")
	}

	if len(f.name) == 0 {
		return nil, errors.New("Field name is empty")
	}

	if _, ok := t.fieldsMap[f.name]; ok {
		return nil, fmt.Errorf("Field already in table: %s", f.name)
	}

	if f.pk && t.pk != nil {
		return nil, fmt.Errorf("PK collision: already assigned to field %s", t.pk.name)
	}
	if f.pk {
		t.pk = f
	}

	f.table = t

	t.fields = append(t.fields, f)
	t.fieldsMap[f.name] = f

	return f, nil
}

func (t *Table) addForeignKey(field *Field, foreignTable *Table, foreignKey *Field) error {
	t.foreignKeys = append(t.foreignKeys,
		&ForeignKey{
			tbl:          t,
			field:        field,
			foreignTable: foreignTable,
			foreignKey:   foreignKey,
		})
	return nil
}

func (t *Table) AddIndex(unique bool, fields ...string) error {
	index := new(Index)
	index.table = t
	index.fields = make([]*Field, len(fields))

	for i := 0; i < len(fields); i++ {
		field := t.Field(fields[i])
		if field == nil {
			return fmt.Errorf("Table.AddIndex: field %s does not exist in table %s", fields[i], t.name)
		}
		index.fields[i] = field
	}
	t.indexes = append(t.indexes, index)
	return nil
}

func (t *Table) CreateTableIndexesSql() ([]string, error) {
	ixs := []string{}
	if len(t.indexes) == 0 {
		return ixs, nil
	}

	for i := 0; i < len(t.indexes); i++ {
		s, err := t.indexes[i].CreateSql(i)
		if err != nil {
			return []string{}, err
		}
		ixs = append(ixs, s)
	}
	return ixs, nil
}

func (idx *Index) CreateSql(n int) (string, error) {
	s := "CREATE "
	if idx.unique {
		s += "UNIQUE "
	}
	s += "INDEX IF NOT EXISTS idx_" + idx.table.name + "_" + strconv.Itoa(n) + " ON " + idx.table.name + "("

	for i := 0; i < len(idx.fields); i++ {
		if i != 0 {
			s += COMMA_SPACE
		}
		s += idx.fields[i].name
	}
	s += ")"

	return s, nil

}

// func (t *Table) SetPrimaryKeyByIndex(i int) error {
// 	if i > len(t.fields) {
// 		return errors.New("Index greater than length of fields")
// 	}
// 	t.pk = t.fields[i]
// 	return nil
// }

// func (t *Table) Validate() error {
// 	return nil
// }

func (t *Table) DropTableSql() (string, error) {
	return "DROP TABLE IF EXISTS " + t.name, nil
}

func (t *Table) CreateTableSql() (string, error) {
	err := checkTable(t)
	if err != nil {
		return "", err
	}

	s := "CREATE TABLE IF NOT EXISTS " + t.name + " ("

	for i := 0; i < len(t.fields); i++ {
		if i != 0 {
			s += COMMA_SPACE
		}
		fs, err := t.fields[i].CreateFieldSql()
		if err != nil {
			return "", err
		}
		s += fs
	}

	s += ")"

	return s, nil
}

func checkTable(t *Table) error {
	if len(t.name) == 0 {
		return errors.New("Table name is empty")
	}

	if len(t.fields) == 0 {
		return errors.New("No fields defined")
	}
	if countPrimaryFields(t.fields) > 1 {
		return errors.New("Multiple primary keys defined")
	}
	return nil
}

func countPrimaryFields(fs []*Field) int {
	n := 0

	for i := 0; i < len(fs); i++ {
		if fs[i].pk {
			n++
		}
	}
	return n
}

func (t *Table) Count(db *sql.DB) (int64, error) {
	row := db.QueryRow("SELECT count(*) from " + t.name)
	var n int64
	err := row.Scan(&n)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return -1, err
	}
	return n, nil
}

func (t *Table) GetMaxId(db *sql.DB) (int64, error) {
	n, err := t.Count(db)
	if err != nil {
		log.Println(err)
		return -1, err
	}

	if n == 0 {
		return 0, nil
	}

	row := db.QueryRow("SELECT max(" + t.pk.name + ") from " + t.name)
	var id int64
	err = row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return -1, err
	}
	return id, nil
}

// FIXXX add bool,
func FindFieldValueById[V *int64 | *float64 | *string](t *Table, db *sql.DB, idValue int64, field *Field, fieldValue V) (bool, error) {
	log.Println(field.name)
	log.Println(t.name)
	log.Println(t.pk)
	log.Println(t.pk.name)
	query := "SELECT " + field.name + " from " + t.name + " where " + t.pk.name + "=?"

	row := db.QueryRow(query, idValue)

	err := row.Scan(fieldValue)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (t *Table) Field(s string) *Field {
	field, _ := t.fieldsMap[s]
	return field
}

func (t *Table) AllFields() string {
	fs := ""
	for i := 0; i < len(t.fields); i++ {
		if i != 0 {
			fs += COMMA_SPACE
		}
		fs += t.fields[i].name
	}
	return fs
}

func (t *Table) SelectRecordsSimpleWhere(db *sql.DB, left, operator, right string, limit, offset int64) (*[]Record, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (t *Table) notField(field string) bool {
	//log.Fatal("Not implemented")
	return false
}

// func (t *Table) SelectRecordSimpleWhere(db *sql.DB, left, operator, right string) (*Record, error) {

// 	if t.notField(left) {
// 		return nil, fmt.Errorf("%s is not a field in table %s", left, t.name)
// 	}

// 	q := "SELECT " + t.AllFields() + " FROM " + t.name + " WHERE " + left + operator + "?"

// 	row := db.QueryRow(q, right)

// 	rec, err := makeRecord(t, row)

// 	if err != nil {
// 		return nil, err
// 	}

// 	return rec, nil
// }

func (t *Table) FieldValueExists(db *sql.DB, field *Field, value any) (bool, error) {
	row := db.QueryRow("SELECT "+t.pk.name+" from "+t.name+" where "+field.name+"= ?",
		value)
	var pk int64
	err := row.Scan(&pk)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
