package dalkeeth

import (
	"database/sql"
)

type Dialect interface {
	CreateTableIndexSql(*Index) (string, error)
	CreateTableSql(*Table) (string, error)
	DialectName() string
	GetSingleRecordSql(*Record, int64) (string, error)
	SaveSql(*Record) (string, error)
	Table(*Table) (string, error)
	ValidTableName(string) error
	JoinSql(*Join, string, ...*Field) error
	ExtractTable(db *sql.DB, tableName string) (*Table, error)
	SelectQuerySql(*SelectQuery) (string, error)
	ArbitraryFunc(string, []any) (string, error)
	FunctionFieldSql(FunctionField) string
	//FieldFunction(int, ...Field)
	//DropTableSql(string)
}
