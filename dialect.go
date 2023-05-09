package dalkeeth

import (
	"database/sql"
)

type Dialect interface {
	ArbitraryFunc(string, []any) (string, error)
	CreateTableIndexSql(*Index) (string, error)
	CreateTableSql(*Table) (string, error)
	DeleteSql(tbl *Table, id int64) (string, error)
	DialectName() string
	ExtractTable(db *sql.DB, tableName string) (*Table, error)
	FieldAsSql(fa *FieldAs) (string, error)
	FunctionFieldSql(FunctionField) (string, error)
	GetSingleRecordSql(*InRecord, int64) (string, error)
	JoinSql(*Join, string, ...*Field) error
	SaveSql(*InRecord) (string, error)
	SelectQuerySql(*SelectQuery) (string, error)
	Table(*Table) (string, error)
	ValidTableName(string) error

	//FieldFunction(int, ...Field)
	//DropTableSql(string)
}
