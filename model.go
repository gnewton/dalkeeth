package dalkeeth

import (
	"fmt"
)

type Model struct {
	tables        []*Table
	tablesMap     map[string]*Table
	selectLimits  map[*Table]int64  // This should be in manager?
	fieldTableMap map[string]*Field // "key=tablename.fieldname", value=*Field
}

func NewModel() *Model {
	m := new(Model)
	m.tablesMap = make(map[string]*Table)
	return m
}
func (m *Model) Table(key string) *Table {
	if t, ok := m.tablesMap[key]; ok {
		return t
	}
	return nil
}

// Key is mneumonic for table; Does not have to be the same as the table sql name.
func (m *Model) AddTable(key string, tbl *Table) error {
	if key == "" {
		return fmt.Errorf("Key is empty string")
	}

	if tbl == nil {
		return fmt.Errorf("Table is nil")
	}

	if t, ok := m.tablesMap[key]; ok {
		return fmt.Errorf("Key %s already occupied by table with name %s", key, t.name)
	}
	m.tablesMap[key] = tbl
	m.tables = append(m.tables, tbl)
	return nil
}
