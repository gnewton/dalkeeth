package dalkeeth

type FunctionField interface {
	AField
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
