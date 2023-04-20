package dalkeeth

type SQLFunctionId int

const (
	//Aggregate functions - https://www.sqlite.org/lang_aggfunc.html#aggfunclist
	AVG SQLFunctionId = iota
	// count(*)
	COUNT //(X)
	// group_concat(X)
	// group_concat(X,Y)
	// max(X)
	// min(X)
	// sum(X)
	// total(X)

	// 	// Scalar functions - https://www.sqlite.org/lang_corefunc.html
	// abs(X)
	// changes()
	// char(X1,X2,...,XN)
	COALESCE //(X,Y,...)
	// format(FORMAT,...)
	// glob(X,Y)
	// hex(X)
	// ifnull(X,Y)
	// iif(X,Y,Z)
	// instr(X,Y)
	// last_insert_rowid()
	// length(X)
	// like(X,Y)
	// like(X,Y,Z)
	// likelihood(X,Y)
	// likely(X)
	// load_extension(X)
	// load_extension(X,Y)
	// lower(X)
	// ltrim(X)
	// ltrim(X,Y)
	MAX //(X,Y,...)
	// min(X,Y,...)
	// nullif(X,Y)
	// printf(FORMAT,...)
	// quote(X)
	RANDOM //()
	// randomblob(N)
	// replace(X,Y,Z)
	// round(X)
	// round(X,Y)
	// rtrim(X)
	// rtrim(X,Y)
	// sign(X)
	// soundex(X)
	// sqlite_compileoption_get(N)
	// sqlite_compileoption_used(X)
	// sqlite_offset(X)
	// sqlite_source_id()
	// sqlite_version()
	// substr(X,Y)
	// substr(X,Y,Z)
	// substring(X,Y)
	// substring(X,Y,Z)
	// total_changes()
	// trim(X)
	// trim(X,Y)
	// typeof(X)
	// unhex(X)
	// unhex(X,Y)
	// unicode(X)
	// unlikely(X)
	// upper(X)
	// zeroblob(N)

	// 	// Time; no arguments
	// 	CURRENT_DATE
	// 	CURRENT_TIME
	// 	DAY
	// 	MONTH
	// 	NOW
	// 	YEAR
	// 	// string functions
	// 	LTRIM
	// 	REPLACE
	// 	RTRIM
	// 	SUBSTRING
	// 	TRIM
	// 	//
	// 	COALESCE
	// 	// MATH functions - https://www.sqlite.org/lang_mathfunc.html
	// acos(X)
	// acosh(X)
	// asin(X)
	// asinh(X)
	// atan(X)
	// atan2(Y,X)
	// atanh(X)
	// ceil(X)
	// ceiling(X)
	// cos(X)
	// cosh(X)
	// degrees(X)
	// exp(X)
	// floor(X)
	// ln(X)
	// log(B,X)
	// log(X)
	// log10(X)
	// log2(X)
	// mod(X,Y)
	// pi()
	// pow(X,Y)
	// power(X,Y)
	// radians(X)
	// sin(X)
	// sinh(X)
	// sqrt(X)
	// tan(X)
	// tanh(X)
	// trunc(X)

)

var sqlFunctionNArgs = map[SQLFunctionId]int{
	AVG:      1,
	COUNT:    3,
	COALESCE: 1,
	MAX:      99,
	RANDOM:   0,
}

type FunctionField struct {
	sqlFunctionId SQLFunctionId
	fields        []AField
}

func (ff FunctionField) ToSqlString(d Dialect) string {
	return d.FunctionFieldSql(ff)
}

func NewFunctionField(sf SQLFunctionId, fields ...AField) AField {
	ff := FunctionField{
		sqlFunctionId: sf,
		fields:        fields,
	}
	return ff
}

func funcToString(name string, v ...any) string {
	return "NotImplemented"
}

type ArbitraryFunc struct {
	name   string
	values []any
}

func (af *ArbitraryFunc) ToSqlString(d Dialect) (string, error) {
	return d.ArbitraryFunc(af.name, af.values)
}
