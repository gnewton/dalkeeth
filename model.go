package dalkeeth

import (
	"fmt"
	"log"
)

type Model struct {
	tables            []*Table
	tablesMap         map[string]*Table
	containsTablesMap map[*Table]any
	selectLimits      map[*Table]int64  // This should be in manager?
	fieldTableMap     map[string]*Field // "key=tablename.fieldname", value=*Field
	frozen            bool
}

func NewModel() *Model {
	m := new(Model)
	m.tablesMap = make(map[string]*Table)
	m.containsTablesMap = make(map[*Table]any)
	return m
}

func (m *Model) Freeze() error {
	if m.frozen {
		return fmt.Errorf("Model is already frozen: multiple freezes?")
	}
	m.fieldTableMap = make(map[string]*Field, 0)

	for i := 0; i < len(m.tables); i++ {
		tbl := m.tables[i]
		// fields
		for i := 0; i < len(tbl.fields); i++ {
			m.fieldTableMap[tbl.name+"."+tbl.fields[i].name] = tbl.fields[i]
			log.Println("Adding table.field", tbl.name+"."+tbl.fields[i].name)
		}
	}
	m.frozen = true
	return nil
}

func (m *Model) HasTable(tbl *Table) bool {
	_, has := m.containsTablesMap[tbl]
	return has
}

func (m *Model) TableByKey(key string) *Table {

	t, ok := m.tablesMap[key]
	if !ok {
		return nil
	}
	return t
	//}
	//return nil
	//return m.tablesMap[key]
}

// Key is mneumonic for table; Does not have to be the same as the table sql name.
func (m *Model) AddTable(key string, tbl *Table) error {
	if m.frozen {
		return fmt.Errorf("Model is frozen: cannot add table to field")
	}
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
	m.containsTablesMap[tbl] = struct{}{}

	return nil
}

func (mdl *Model) AddForeignKey(tbl *Table, field string, foreignTbl *Table, foreignKeyField string) error {
	if mdl.frozen {
		return fmt.Errorf("Model is frozen: change")
	}
	if tbl == nil {
		return fmt.Errorf("manager.AddForeignKey: Table is nil")
	}

	if foreignTbl == nil {
		return fmt.Errorf("manager.AddForeignKey: Foreign table is nil")
	}

	if field == "" {
		return fmt.Errorf("manager.AddForeignKey: Field name is empty")
	}
	if foreignKeyField == "" {
		return fmt.Errorf("manager.AddForeignKey: foreignKeyField is empty")
	}

	if tbl.Field(field) == nil {
		return fmt.Errorf("manager.AddForeignKey: Field %s does not exist in table %s", field, tbl.name)
	}

	if foreignTbl.Field(foreignKeyField) == nil {
		return fmt.Errorf("manager.AddForeignKey: Foreign key field %s does not exist in table %s", foreignKeyField, foreignTbl.name)
	}

	f, ok := tbl.fieldsMap[field]
	if !ok {
		return fmt.Errorf("Field %s not found in table %s", field, tbl.name)
	}

	fk, ok := foreignTbl.fieldsMap[foreignKeyField]
	if !ok {
		return fmt.Errorf("Field %s not found in table %s", foreignKeyField, foreignTbl.name)
	}

	return tbl.addForeignKey(f, foreignTbl, fk)
}
