package dalkeeth

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Manager struct {
	model *Model
	db    *sql.DB
	tx    *sql.Tx
	//tables []*Table
	//tablesMap     map[string]*Table
	selectLimits  map[*Table]int64
	dialect       Dialect
	fieldTableMap map[string]*Field // "key=tablename.fieldname", value=*Field
}

func NewManagerWithModel(model *Model) *Manager {
	m := new(Manager)
	m.model = model
	m.selectLimits = make(map[*Table]int64)
	return m
}

func NewManager() *Manager {
	model := NewModel()
	return NewManagerWithModel(model)
}

func (m *Manager) Close() error {
	if m.db == nil {
		return fmt.Errorf("Trying to close nil db")
	}
	return m.db.Close()
}

func (m *Manager) Table(key string) *Table {
	return m.model.Table(key)
}

// Key is mneumonic for table; Does not have to be the same as the table sql name.

func (m *Manager) CreateTablesSQL() ([]string, error) {
	if m.dialect == nil {
		return nil, fmt.Errorf("Dialect is nil")
	}
	var sql []string

	for i := 0; i < len(m.model.tables); i++ {
		s, err := m.dialect.CreateTableSql(m.model.tables[i])
		if err != nil {
			return nil, err
		}
		sql = append(sql, s)
	}

	return sql, nil
}

func (m *Manager) CreateTableIndexesSQL() ([]string, error) {
	if m.dialect == nil {
		return nil, fmt.Errorf("Dialect is nil")
	}
	var sql []string

	for i := 0; i < len(m.model.tables); i++ {
		for j := 0; j < len(m.model.tables[i].indexes); j++ {
			s, err := m.dialect.CreateTableIndexSql(m.model.tables[i].indexes[j])
			if err != nil {
				return nil, err
			}
			sql = append(sql, s)
		}
	}

	return sql, nil
}

func (m *Manager) BatchChannel(chunkSize int) (chan *Record, error) {
	return nil, NotImplemented
}

func (m *Manager) Batch(recs []*Record) error {
	var err error
	if m.tx == nil {
		return errors.New("manager.Save: tx is nil")
	}

	saveSql, err := m.dialect.SaveSql(recs[0])
	if err != nil {
		return err
	}

	stmt, err := m.tx.Prepare(saveSql)
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

func (m *Manager) save(r *Record) error {
	return NotImplemented
}

func (m *Manager) Update(r *Record) error {
	return NotImplemented
}

func (m *Manager) GetNamed(tableName string, id int64) (*Record, error) {
	// find tableName
	// Get(tablename, id)
	return nil, NotImplemented
}

func (m *Manager) Get(tbl *Table, id int64) (*Record, error) {
	if id < 0 {
		return nil, errors.New("manager.Get: id < 0: ")
	}
	if m.db == nil {
		return nil, errors.New("manager.Get: db is nil")
	}
	if m.dialect == nil {
		return nil, errors.New("manager.Save: dialext is nil")
	}

	rec := tbl.NewRecord()

	query, err := m.dialect.GetSingleRecordSql(rec, id)
	if err != nil {
		return nil, err
	}
	log.Println(query)

	values, err := rawWantedValues(rec.values)
	if err != nil {
		return nil, err
	}
	log.Println("-------------------rawWantedValues", values)

	row := m.db.QueryRow(query)
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
		log.Println("Manager", i)
		//}
	}

	return rec, nil

}

// Using db, not tx
func (m *Manager) Save(r *Record) error {
	if r == nil {
		return errors.New("manager.Save: record is nil")
	}
	if m.dialect == nil {
		return errors.New("manager.Save: dialext is nil")
	}
	saveSql, err := m.dialect.SaveSql(r)
	if err != nil {
		return err
	}

	log.Println(saveSql)

	rawValues := rawValues(r.values)
	_, err = m.db.Exec(saveSql, rawValues...)

	if err != nil {
		return err
	}
	return nil
}

func (m *Manager) SaveTx(r *Record) error {
	if m.tx == nil {
		return errors.New("manager.Save: tx is nil")
	}
	if m.dialect == nil {
		return errors.New("manager.Save: dialext is nil")
	}
	saveSql, err := m.dialect.SaveSql(r)
	if err != nil {
		return err
	}

	log.Println(saveSql)

	rawValues := rawValues(r.values)
	_, err = m.tx.Exec(saveSql, rawValues...)

	if err != nil {
		log.Println(saveSql)
		log.Println(err)
		return err
	}
	return nil
}

func (m *Manager) Begin() error {
	if m.db == nil {
		return errors.New("DB is nil")
	}

	if m.tx != nil {
		return errors.New("Already in transaction")
	}
	var err error

	m.tx, err = m.db.Begin()

	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) Commit() error {
	if m.tx == nil {
		return errors.New("No transaction started")
	}

	defer func() {
		m.tx = nil
	}()

	err := m.tx.Commit()
	if err != nil {
		log.Println("Commit error")
		log.Println(err)
		rollbackErr := m.tx.Rollback()
		if rollbackErr != nil {
			log.Println("Rollback error")
			log.Println(rollbackErr)
			return rollbackErr
		}
		return err
	}
	return nil
}

func (m *Manager) Rollback() error {
	return nil
}
