package dalkeeth

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type DialectSqlite3 struct {
}

func (d *DialectSqlite3) DialectName() string {
	return "Sqlite3"
}

func (d *DialectSqlite3) Table(t *Table) (string, error) {
	if err := d.ValidTableName(t.name); err != nil {
		return "", err
	}
	return "[" + t.name + "]", nil
}

func (d *DialectSqlite3) CreateTableIndexSql(ind *Index) (string, error) {
	if ind == nil {
		return "", errors.New("Index is nil")
	}

	s := "CREATE "
	if ind.unique {
		s += "UNIQUE"
	}
	s += " INDEX idx_" + ind.table.name

	for i := 0; i < len(ind.fields); i++ {
		s += "_" + ind.fields[i].name
	}

	s += " ON " + ind.table.name + "("

	for i := 0; i < len(ind.fields); i++ {
		if i != 0 {
			s += COMMA_SPACE
		}
		s += ind.fields[i].name
	}

	s += ")"

	return s, nil
}
func (d *DialectSqlite3) SelectQuerySql(q *SelectQuery) (string, error) {

	if err := q.Validate(); err != nil {
		return "", err
	}

	s := "SELECT "
	if q.distinct {
		s += "DISTINCT "
	}

	s += d.makeFields(q.fields)
	if len(q.pks) != 0 {
		s += d.makePks(q.pks)
	} else {
		s += d.makeWhereClause(q.where)
	}
	s += d.makeGroupBy(q.groupBy)
	s += d.makeHaving(q.having)
	if q.limit > 0 {
		s += " LIMIT " + strconv.FormatInt(q.limit, 10)
	}
	if q.offset > 0 {
		s += " OFFSET " + strconv.FormatInt(q.offset, 10)
	}

	if q.ordering != NoOrdering {
		s += "ORDER BY "
		s += d.makeOrderByFields(q.orderBy)
		if q.ordering == ASC {
			s += " ASC"
		} else {
			s += " DESC"
		}
	}

	log.Fatal("NotImplemented")

	return "", nil
}

func (d *DialectSqlite3) CreateTableSql(t *Table) (string, error) {
	var err error
	if t == nil {
		return "", errors.New("Table is nil")
	}

	if err = d.ValidTableName(t.name); err != nil {
		return "", err
	}

	s := "CREATE TABLE " + t.name

	if len(t.fields) == 0 {
		return "", errors.New("Num fields = zero")
	}

	s += " ("

	var fsql string
	for i := 0; i < len(t.fields); i++ {
		if fsql, err = d.fieldSql(t.fields[i]); err != nil {
			return "", err
		}
		if i != 0 {
			s += COMMA_SPACE
		}
		s += fsql
	}

	foreignKeysSql, err := d.foreignKeys(t.foreignKeys)
	if err != nil {
		return "", err
	}
	s += foreignKeysSql
	s += ")"

	return s, nil
}

func (d *DialectSqlite3) foreignKeys(fKeys []*ForeignKey) (string, error) {
	if len(fKeys) == 0 {
		return "", nil
	}

	var s string
	for i := 0; i < len(fKeys); i++ {
		fk := fKeys[i]
		s += ", FOREIGN KEY(" + fk.field.name + ") REFERENCES " + fk.foreignTable.name + "(" + fk.field.name + ")"
	}

	return s, nil
}

func (d *DialectSqlite3) fieldSql(f *Field) (string, error) {
	s := f.name + SPACE

	switch f.fieldType {
	case IntType:
		s += "INT"
	case StringType:
		if f.length == 0 {
			s += "TEXT"
		} else {
			s += "varchar(" + strconv.Itoa(f.length) + ")"
		}
	case FloatType:
		s += "REAL"

	case BoolType:
		s += "BOOLEAN"
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
	if f.defaultValue != "" {
		defaultValue, err := makeDefault(f)
		if err != nil {
			return "", err
		}
		s += defaultValue
	}
	return s, nil
}

// Verify the string of the default type is the type of the field; i.e. "43.5" is float; "231" is int;
// No need to check string:string
func defaultTypeMatchesFieldType(f *Field) error {
	var err error
	switch f.fieldType {
	case IntType:
		_, err = strconv.ParseUint(f.defaultValue, 10, 64)
	case FloatType:
		_, err = strconv.ParseFloat(f.defaultValue, 64)
	}
	return err
}

func makeDefault(f *Field) (string, error) {
	if err := defaultTypeMatchesFieldType(f); err != nil {
		return "", err
	}

	def := " DEFAULT "
	//return `" + quote(f.defaultValue) + "`"
	var delimiter = ""
	if f.fieldType == StringType {
		delimiter = "`"
	}
	def += delimiter + quote(f.defaultValue) + delimiter

	return def, nil
}

func quote(s string) string {
	return strings.ReplaceAll(s, "`", "``")
}

func (d *DialectSqlite3) ValidTableName(name string) error {
	if len(name) == 0 || strings.HasPrefix(name, "sqlite_") {
		return fmt.Errorf("Invalid table name [%s] for dialect %s", name, d.DialectName())
	}

	return nil

}

func closeRows(rows *sql.Rows) {
	if rows != nil {
		err := rows.Close()
		if err != nil {
			log.Println(err)
		}
	}
}

func (d *DialectSqlite3) ExtractTable(db *sql.DB, tableName string) (*Table, error) {
	return d.tableInfo(db, tableName)
}

func (d *DialectSqlite3) tableInfo(db *sql.DB, tableName string) (*Table, error) {
	q := fmt.Sprintf("PRAGMA table_info(%s);", tableName)
	var cid, notNull, pk int64
	var name, ftype string
	var dflt_value sql.NullString

	log.Println(q)
	rows, err := db.Query(q)
	defer closeRows(rows)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		if err := rows.Scan(&cid, &name, &ftype, &notNull, &dflt_value, &pk); err != nil {
			return nil, err
		}
		log.Println("tableName:", cid, name, ftype, notNull, "[", dflt_value, "]", pk)
	}

	return nil, nil
}

// Only use fields that have set values
func setRecords(r *Record) (string, error) {
	s := ""

	first := true

	for i := 0; i < len(r.values); i++ {
		v := r.values[i]

		if v.isSet {
			if !first {
				s += COMMA_SPACE
			}
			first = false
			s += v.field.name
		} else {
			if v.field.notNull {
				return "", errors.New("Field " + v.field.name + " must be set: not null")
			}
			if v.field.pk {
				return "", errors.New("Field " + v.field.name + " must be set: primary key")
			}
		}
	}
	if first {
		return "", errors.New("No fields set in record")
	}
	return s, nil
}

// Only use fields that have set values
func wantedFields(r *Record) (string, error) {
	s := ""

	first := true

	for i := 0; i < len(r.values); i++ {
		v := r.values[i]

		if v.isWanted {
			if !first {
				s += COMMA_SPACE
			}
			first = false
			s += v.field.name
		}
	}
	if first {
		return "", errors.New("No fields set in record")
	}
	return s, nil
}

func (d *DialectSqlite3) GetSingleRecordSql(rec *Record, id int64) (string, error) {

	if rec == nil {
		return "", errors.New("Record is nil")
	}

	s := "SELECT "

	wanted, err := wantedFields(rec)
	if err != nil {
		return "", err
	}

	s += wanted + " FROM " + rec.table.name + " WHERE " + ID + "=" + strconv.FormatInt(id, 10)

	return s, nil
}

func (d *DialectSqlite3) valuesPlaceholders(r *Record) (string, error) {
	s := ""

	first := true
	placeHolderCount := 1
	for i := 0; i < len(r.values); i++ {
		v := r.values[i]

		if v.isSet {
			if !first {
				s += COMMA_SPACE
			}
			first = false
			s += "$" + strconv.Itoa(placeHolderCount)
			placeHolderCount++
		}
	}
	if first {
		return "", errors.New("DialectSqlite3.valuesPlaceholders: No fields set in record")
	}
	return s, nil
}

func (d *DialectSqlite3) SaveSql(r *Record) (string, error) {
	if r == nil {
		return "", errors.New("DialectSqlite3.Save: record is nil")
	}

	s := "INSERT INTO " + r.table.name + " ("

	fieldsSet, err := setRecords(r)
	if err != nil {
		return "", err
	}

	s += fieldsSet
	s += ") VALUES ("

	valuesP, err := d.valuesPlaceholders(r)
	if err != nil {
		return "", err
	}

	s += valuesP
	s += ")"

	return s, nil
}

func (d *DialectSqlite3) JoinSql(*Join, string, ...*Field) error {
	return NotImplemented
}

func (d *DialectSqlite3) makeFields(fields []*SelectField) string {
	if len(fields) == 0 {
		return ""
	}

	s := fields[0].name

	for i := 1; i < len(fields); i++ {
		s += COMMA_SPACE + fields[i].name
	}

	return s
}

func (d *DialectSqlite3) makePks(pks []int64) string {
	return "fail"
}

func (d *DialectSqlite3) makeWhereClause(cond Condition) string {
	return "fail"
}
func (d *DialectSqlite3) makeGroupBy(fields []*Field) string {
	if len(fields) == 0 {
		return ""
	}

	s := "GROUP BY " + fields[0].name
	for i := 1; i < len(fields); i++ {
		s += COMMA_SPACE + fields[i].name
	}
	return s
}

func (d *DialectSqlite3) makeHaving(cond Condition) string {
	return "fail"
}

const SPACE = " "
const COMMA_SPACE = ", "

func (d *DialectSqlite3) makeOrderByFields(fields []*SelectField) string {
	if len(fields) == 0 {
		return ""
	}

	var s string

	for i := 0; i < len(fields); i++ {
		if i != 0 {
			s += COMMA_SPACE
		}
		f := fields[i]
		s += f.name
		if f.ordering != NoOrdering {
			s += SPACE
		}
		s += f.ordering.String()
	}
	return s
}

func (d *DialectSqlite3) ArbitraryFunc(string, []any) (string, error) {
	return "", NotImplemented
}

func (d *DialectSqlite3) FunctionFieldSql(ff FunctionField) string {
	switch ff.sqlFunction{
		case 
	}
	log.Fatal(NotImplemented)
	return ""
}
