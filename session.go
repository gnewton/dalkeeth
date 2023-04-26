package dalkeeth

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Session struct {
	model *Model
	db    *sql.DB
	tx    *sql.Tx
	//tables []*Table
	//tablesMap     map[string]*Table
	selectLimits  map[*Table]int64
	dialect       Dialect
	fieldTableMap map[string]*Field // "key=tablename.fieldname", value=*Field
}

func NewSession(model *Model) (*Session, error) {
	if !model.frozen {
		return nil, fmt.Errorf("Model needs to be frozen before using it")
	}
	m := new(Session)
	m.model = model
	m.selectLimits = make(map[*Table]int64)
	return m, nil
}

func (sess *Session) Close() error {
	if sess.db == nil {
		return fmt.Errorf("Trying to close nil db")
	}
	return sess.db.Close()
}

func (sess *Session) Table(key string) (*Table, bool) {
	return sess.model.Table(key)
}

func (sess *Session) InstantiateModel() error {
	return nil
}

func (sess *Session) createTablesSQL() ([]string, error) {
	if sess.dialect == nil {
		return nil, fmt.Errorf("Dialect is nil")
	}
	var sql []string

	for i := 0; i < len(sess.model.tables); i++ {
		s, err := sess.dialect.CreateTableSql(sess.model.tables[i])
		if err != nil {
			return nil, err
		}
		sql = append(sql, s)
	}

	return sql, nil
}

func (sess *Session) createTableIndexesSQL() ([]string, error) {
	if sess.dialect == nil {
		return nil, fmt.Errorf("Dialect is nil")
	}
	var sql []string

	for i := 0; i < len(sess.model.tables); i++ {
		for j := 0; j < len(sess.model.tables[i].indexes); j++ {
			s, err := sess.dialect.CreateTableIndexSql(sess.model.tables[i].indexes[j])
			if err != nil {
				return nil, err
			}
			sql = append(sql, s)
		}
	}

	return sql, nil
}

func (sess *Session) BatchChannel(chunkSize int) (chan *Record, error) {
	return nil, NotImplemented
}

func (sess *Session) Batch(recs []*Record) error {
	var err error
	if sess.tx == nil {
		return errors.New("manager.Save: tx is nil")
	}

	saveSql, err := sess.dialect.SaveSql(recs[0])
	if err != nil {
		return err
	}

	stmt, err := sess.tx.Prepare(saveSql)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // danger!

	for i := 0; i < len(recs); i++ {
		rawValues := rawValues(recs[i].values)
		_, err = stmt.Exec(rawValues...)
		if err != nil {
			log.Println("manager.Batch: error")
			log.Println(err)
			return err
		}
	}
	return nil
}

func (sess *Session) save(r *Record) error {
	return NotImplemented
}

func (sess *Session) Update(r *Record) error {
	return NotImplemented
}

func (sess *Session) GetNamed(tableName string, id int64) (*Record, error) {
	// find tableName
	// Get(tablename, id)
	return nil, NotImplemented
}

func (sess *Session) GetS(tblName string, id int64) (*Record, error) {
	var tbl *Table
	var ok bool

	if tbl, ok = sess.model.tablesMap[tblName]; !ok {
		return nil, fmt.Errorf("Unknown table name:[%s]", tblName)
	}

	return sess.Get(tbl, id)
}

func (sess *Session) Get(tbl *Table, id int64) (*Record, error) {
	if id < 0 {
		return nil, errors.New("manager.Get: id < 0: ")
	}
	if sess.db == nil {
		return nil, errors.New("manager.Get: db is nil")
	}
	if sess.dialect == nil {
		return nil, errors.New("manager.Save: dialext is nil")
	}

	rec := tbl.NewRecord()

	query, err := sess.dialect.GetSingleRecordSql(rec, id)
	if err != nil {
		return nil, err
	}
	log.Println(query)

	values, err := rawWantedValues(rec.values)
	if err != nil {
		return nil, err
	}
	log.Println("-------------------rawWantedValues", values)

	row := sess.db.QueryRow(query)
	err = row.Scan(values...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	for i := 0; i < len(rec.values); i++ {
		v := rec.values[i]
		//if v.isWanted {
		actual(v)
		log.Println("Session", i)
		//}
	}

	return rec, nil

}

// Using db, not tx
func (sess *Session) Save(r *Record) error {
	if r == nil {
		return errors.New("manager.Save: record is nil")
	}
	if sess.dialect == nil {
		return errors.New("manager.Save: dialext is nil")
	}
	saveSql, err := sess.dialect.SaveSql(r)
	if err != nil {
		return err
	}

	log.Println(saveSql)

	rawValues := rawValues(r.values)
	_, err = sess.db.Exec(saveSql, rawValues...)

	if err != nil {
		return err
	}
	return nil
}

func (sess *Session) SaveTx(r *Record) error {
	if sess.tx == nil {
		return errors.New("manager.Save: tx is nil")
	}
	if sess.dialect == nil {
		return errors.New("manager.Save: dialext is nil")
	}
	saveSql, err := sess.dialect.SaveSql(r)
	if err != nil {
		return err
	}

	log.Println(saveSql)

	rawValues := rawValues(r.values)
	_, err = sess.tx.Exec(saveSql, rawValues...)

	if err != nil {
		log.Println(saveSql)
		log.Println(err)
		return err
	}
	return nil
}

func (sess *Session) Begin() error {
	if sess.db == nil {
		return errors.New("DB is nil")
	}

	if sess.tx != nil {
		return errors.New("Already in transaction")
	}
	var err error

	sess.tx, err = sess.db.Begin()

	if err != nil {
		return err
	}

	return nil
}

func (sess *Session) Commit() error {
	if sess.tx == nil {
		return errors.New("No transaction started")
	}

	defer func() {
		sess.tx = nil
	}()

	err := sess.tx.Commit()
	if err != nil {
		log.Println("Commit error")
		log.Println(err)
		rollbackErr := sess.tx.Rollback()
		if rollbackErr != nil {
			log.Println("Rollback error")
			log.Println(rollbackErr)
			return rollbackErr
		}
		return err
	}
	return nil
}

func (sess *Session) Rollback() error {
	return nil
}

func (sess *Session) NewSelectQuery() *SelectQuery {
	q := SelectQuery{
		Fields: make([]AField, 0),
	}
	return &q
}

func (sess *Session) Exists(t *Table, id int64) (bool, error) {
	return recordExists(sess.db, t.name, id)

}
