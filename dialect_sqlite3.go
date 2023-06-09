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

	if !q.Validated() {
		return "", errors.New("SelectQuery not validated")
	}

	s := "SELECT "
	if q.Distinct {
		s += "DISTINCT "
	}

	s += d.makeFields(q.Fields)
	if len(q.Pks) != 0 {
		s += d.makePks(q.Pks)
	} else {
		s += d.makeWhereClause(q.Where)
	}
	s += d.makeGroupBy(q.GroupBy)
	s += d.makeHaving(q.Having)
	if q.Limit > 0 {
		s += " LIMIT " + strconv.FormatInt(q.Limit, 10)
	}
	if q.Offset > 0 {
		s += " OFFSET " + strconv.FormatInt(q.Offset, 10)
	}

	if q.GlobalOrdering != NoOrdering {
		s += "ORDER BY "
		s += d.makeOrderByFields(q.OrderByFields)
		if q.GlobalOrdering == ASC {
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

	s := "CREATE TABLE IF NOT EXISTS " + t.name

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
	if db == nil {
		return nil, errors.New("DB is nil")
	}
	q := fmt.Sprintf("PRAGMA table_info(%s);", tableName)
	var cid, notNull, pk int64
	var name, ftype string
	var dflt_value sql.NullString

	//log.Println(q)
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
func insertFields(r *InRecord) (string, error) {
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
func wantedFields(r *InRecord) (string, error) {
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

func (d *DialectSqlite3) GetSingleRecordSql(rec *InRecord, id int64) (string, error) {

	if rec == nil {
		return "", errors.New("InRecord is nil")
	}

	s := "SELECT "

	wanted, err := wantedFields(rec)
	if err != nil {
		return "", err
	}

	s += wanted + " FROM " + rec.table.name + " WHERE " + ID + "=" + strconv.FormatInt(id, 10)

	return s, nil
}

func (d *DialectSqlite3) valuesPlaceholders(r *InRecord) (string, error) {
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

func (d *DialectSqlite3) SaveSql(r *InRecord) (string, error) {
	if r == nil {
		return "", errors.New("DialectSqlite3.Save: record is nil")
	}

	s := "INSERT INTO " + r.table.name + " ("

	fieldsSet, err := insertFields(r)
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

func (d *DialectSqlite3) DeleteSql(tbl *Table, id int64) (string, error) {
	if tbl == nil {
		return "", errors.New("DialectSqlite3.Delete: table is nil")
	}
	if id < 0 {
		return "", errors.New("DialectSqlite3.Delete: id < 0")
	}

	if tbl.name == "" {
		return "", errors.New("DialectSqlite3.Delete: table name is empty string")
	}

	return "DELETE FROM " + tbl.name + " WHERE " + tbl.pk.name + "=?", NotImplemented
}

func (d *DialectSqlite3) JoinSql(*Join, string, ...*Field) error {
	return NotImplemented
}

func (d *DialectSqlite3) makeFields(fields []AField) string {
	if len(fields) == 0 {
		return ""
	}

	//s := fields[0].name
	//s := fields[0].name
	s, _ := fields[0].ToSqlString(d)

	for i := 1; i < len(fields); i++ {
		field, _ := fields[i].ToSqlString(d)
		s += COMMA_SPACE + field
	}

	return s
}

func (d *DialectSqlite3) makePks(pks []int64) string {
	log.Fatal("makePks")
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

func (d *DialectSqlite3) makeOrderByFields(fields []*FieldOrdered) string {
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

var sqlFunctionString = map[SQLFunctionId]string{
	AVG:      "AVG",
	COUNT:    "COUNT",
	COALESCE: "COALESCE",
	MAX:      "MAX",
	RANDOM:   "RANDOM",
}

func (d *DialectSqlite3) FunctionFieldSql(ff FunctionField) (string, error) {
	if true {
		return "HEOOO", nil
	}
	var numExpectedFields int
	var exists bool
	if numExpectedFields, exists = sqlFunctionNArgs[ff.sqlFunctionId]; !exists {
		log.Fatal(fmt.Errorf("FunctionField does not exist: %d   in sqlFunctionId map", ff.sqlFunctionId))
	}
	if numExpectedFields != len(ff.fields) {
		log.Fatal(fmt.Errorf("NumExpectedFields %d for functionid %d does not match acting num fields %d", numExpectedFields, ff.sqlFunctionId, len(ff.fields)))
	}

	var funcName string
	if funcName, exists = sqlFunctionString[ff.sqlFunctionId]; !exists {
		log.Fatal(fmt.Errorf("FunctionField does not exist: %d in sqlFunctionString map", ff.sqlFunctionId))
	}

	f := funcName + "("

	for i := 0; i < len(ff.fields); i++ {
		if i > 0 {
			f += ","
		}
		s, err := ff.fields[i].ToSqlString(d)
		if err != nil {
			return "", err
		}
		f += s
	}
	f += ")"

	return f, nil
}

func (d *DialectSqlite3) FieldAsSql(fa *FieldAs) (string, error) {
	if fa == nil {
		return "", errors.New("FieldAs is nil")
	}

	if fa.field == nil {
		return "", errors.New("FieldAs.field is nil")
	}

	if len(fa.field.name) == 0 {
		return "", errors.New("FieldAs.field.name is empty string")
	}

	return fa.field.name + " AS " + fa.alias, nil
}

func (d *DialectSqlite3) SelectQuerySql2(q *Query) (string, error) {
	var sql string = "SELECT "

	err := makeSelectFields(&sql, q.selectFields, q.selectRaw)
	if err != nil {
		return "", err
	}

	sql += " FROM "

	err = makeFromTables(&sql, q.fromTables, q.fromRaw)
	if err != nil {
		return "", err
	}

	err = makeWhereClause(&sql, q.where, q.whereEquals, q.whereRaw)
	if err != nil {
		return "", err
	}

	return sql, nil
}

func makeSelectFields(sql *string, fields []*Field, rawFields []string) error {
	for i := 0; i < len(fields); i++ {
		if i != 0 {
			*sql += ","
		}
		*sql += fields[i].name
	}

	for i := 0; i < len(rawFields); i++ {
		if i != 0 || len(fields) > 0 {
			*sql += ","
		}
		*sql += rawFields[i]
	}

	return nil
}

func makeFromTables(sql *string, tables []*Table, rawTables string) error {
	for i := 0; i < len(tables); i++ {
		if i != 0 {
			*sql += ","
		}
		*sql += tables[i].name
	}

	if len(tables) > 0 && len(rawTables) > 0 {
		*sql += ","
	}
	*sql += rawTables

	*sql += " "

	return nil
}

func makeWhereClause(sql *string, where []*Condition, whereEquals []AField, whereRaw string) error {

	*sql += "WHERE "

	return nil
}
