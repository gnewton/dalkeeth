package dalkeeth

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Session struct {
	model         *Model
	db            *sql.DB
	tx            *sql.Tx
	selectLimits  map[*Table]int64
	dialect       Dialect
	fieldTableMap map[string]*Field // "key=tablename.fieldname", value=*Field
	readWrite     bool
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

func (sess *Session) TableByKey(key string) *Table {
	return sess.model.TableByKey(key)
}

func (sess *Session) InstantiateModel() error {
	if !sess.readWrite {
		return errors.New("InstantiateModel: session is read-only")
	}
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

// Creates new tx and commits at end
func (sess *Session) BatchChannel(chunkSize int) (chan *Record, error) {
	if !sess.readWrite {
		return nil, errors.New("BatchChannel: session is read-only")
	}
	return nil, NotImplemented
}

// Creates new tx and commits at end
func (sess *Session) Batch(recs []*Record) error {
	if !sess.readWrite {
		return errors.New("Batch: session is read-only")
	}
	var err error
	if sess.tx != nil {
		return errors.New("session.Save: tx is not nil")
	}
	err = sess.Begin()
	if err != nil {
		return err
	}

	saveSql, err := sess.dialect.SaveSql(recs[0])
	if err != nil {
		//FIXXX: roll back
		return err
	}

	stmt, err := sess.tx.Prepare(saveSql)
	if err != nil {
		//FIXXX: roll back
		log.Fatal(err)
	}
	defer stmt.Close() // danger!

	for i := 0; i < len(recs); i++ {
		rawValues := rawValues(recs[i].values)
		_, err = stmt.Exec(rawValues...)
		if err != nil {
			log.Println("session.Batch: error")
			log.Println(err)
			//FIXXX: roll back
			return err
		}
	}
	err = sess.Commit()
	if err != nil {
		//FIXXX: roll back
		return err
	}

	return nil
}

func (sess *Session) Update(r *Record) error {
	if !sess.readWrite {
		return errors.New("Update: session is read-only")
	}
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
		return nil, errors.New("session.Get: id < 0: ")
	}
	if sess.db == nil {
		return nil, errors.New("session.Get: db is nil")
	}
	if sess.dialect == nil {
		return nil, errors.New("session.Save: dialext is nil")
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
	if !sess.readWrite {
		return errors.New("Save: session is read-only")
	}
	if r == nil {
		return errors.New("session.Save: record is nil")
	}
	if len(r.values) == 0 {
		return errors.New("session.Save:record.values is empty")
	}
	if sess.dialect == nil {
		return errors.New("session.Save: dialext is nil")
	}
	if r.table == nil {
		return errors.New("session.Save:record.Table is nil")
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
	if !sess.readWrite {
		return errors.New("SaveTx: session is read-only")
	}
	if sess.tx == nil {
		return errors.New("session.Save: tx is nil")
	}
	if sess.dialect == nil {
		return errors.New("session.Save: dialext is nil")
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

// Using db, not tx
func (sess *Session) Delete(tbl *Table, id int64) error {
	if !sess.readWrite {
		return errors.New("Save: session is read-only")
	}
	if tbl == nil {
		return errors.New("session.Delete: table is nil")
	}
	if id < 0 {
		return errors.New("session.Delete: id < 0")
	}
	if sess.dialect == nil {
		return errors.New("session.Save: dialext is nil")
	}

	deleteSql, err := sess.dialect.DeleteSql(tbl, id)
	if err != nil {
		return err
	}

	log.Println(deleteSql)

	_, err = sess.db.Exec(deleteSql, id)

	if err != nil {
		return err
	}
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
