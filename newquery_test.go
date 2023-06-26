package dalkeeth

import (
	"testing"
)

func TestNewQuery_T1(t *testing.T) {
	setupTest()
	//mgr, err := testModel0()
	_, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}
	// end setup
}

func TestNewQuery_T2(t *testing.T) {
	model, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}

	dialect := new(DialectSqlite3)

	q := NewQuery()
	f0 := model.TableField(TPerson, FId)
	q.selectFields = append(q.selectFields, f0)
	// FIXXX
	//q.rawFields = append(q.selectFields, f0)

	sql, err := dialect.SelectQuerySql2(q)

	t.Log(sql)
	if err != nil {
		t.Fatal(err)
	}

}

func TestNewQuery_RawAll(t *testing.T) {

	// Make model
	model, err := testModel0()
	if err != nil {
		t.Fatal(err)
	}
	// Make session
	sess, err := NewSession(model)
	if err != nil {
		t.Fatal(err)
	}
	// Get SQL dialect
	dialect := new(DialectSqlite3)
	sess.dialect = dialect

	//
	err = sess.WriteModelTableSchemaToDB()
	if err != nil {
		t.Fatal(err)
	}

	q := NewQuery().SelectByName(FId, FName).FromRaw("persons").WhereRaw("name like 'S*'")

	sql, err := dialect.SelectQuerySql2(q)

	t.Log(sql)
	t.Fatal(err)

}
