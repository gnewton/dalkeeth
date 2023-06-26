package dalkeeth

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func writeTestTableRecords(sess *Session) error {

	return NotImplemented
	//	return nil

}

func writeTestModelSchema(mdl *Model) (*Session, error) {
	sess, err := NewSession(mdl)
	if err != nil {
		return nil, err
	}
	sess.dialect = new(DialectSqlite3)
	sqls, err := sess.createTablesSQL()

	if err != nil {
		return nil, err
	}

	indexesSql, err := sess.createTableIndexesSQL()

	if err != nil {
		return nil, err
	}

	db, err := openTestDB()
	if err != nil {
		return nil, err
	}

	sess.db = db
	sess.dialect = new(DialectSqlite3)

	// Create tables sql
	log.Println("Sql tables:", sqls)
	for i := 0; i < len(sqls); i++ {
		createSql := sqls[i]
		log.Println(createSql)
		fmt.Fprintln(os.Stdout, createSql)
		result, err := db.Exec(createSql)

		if err != nil {
			log.Println(fmt.Errorf("writeTestModelSchema: %s", err))
			//return nil, err
			return nil, fmt.Errorf("DB.Exec error: %w", err)
		}
		_, err = result.RowsAffected()
		if err != nil {
			//log.Println(fmt.Errorf("writeTestModelSchema: %s", err))
			return nil, fmt.Errorf("result.RowsAffected: %w", err)
			//return nil, err
		}

	}

	// Create table indexes sql
	for i := 0; i < len(indexesSql); i++ {
		createSql := indexesSql[i]
		result, err := db.Exec(createSql)

		if err != nil {
			return nil, err
		}
		_, err = result.RowsAffected()
		if err != nil {
			return nil, err
		}

	}

	return sess, nil
}

const TClients = "clients"

func fullModel() (*Model, error) {
	mdl := new(Model)
	clients, err := mdl.NewTable(TClients)
	if err != nil {
		return nil, err
	}

	// Client: id, name, city
	clients.AddFields([]*Field{
		&Field{
			name:      "id",
			fieldType: IntType,
			pk:        true,
		},
		&Field{
			name:      "name",
			fieldType: StringType,
			length:    64,
			notNull:   true,
		},
		&Field{
			name:      "city",
			fieldType: IntType,
		},
	}...)

	// Order: id, clientId
	orders, err := mdl.NewTable("orders")
	if err != nil {
		return nil, err
	}

	orders.AddFields([]*Field{
		&Field{
			name:      "clientId",
			fieldType: IntType,
		},
		&Field{
			name:      "orderId",
			fieldType: IntType,
		},
	}...)
	// Product: id, name
	products, err := mdl.NewTable("products")
	if err != nil {
		return nil, err
	}

	products.AddFields([]*Field{
		&Field{
			name:      "id",
			fieldType: IntType,
			pk:        true,
		},
		&Field{
			name:      "name",
			fieldType: StringType,
		},
	}...)

	// OrderProductJoin: clientId, productId
	orderProducts, err := mdl.NewTable("order_products_join")
	if err != nil {
		return nil, err
	}

	orderProducts.AddFields([]*Field{
		&Field{
			name:      "id",
			fieldType: IntType,
			pk:        true,
		},
		&Field{
			name:      "orderId",
			fieldType: IntType,
		},
		&Field{
			name:      "productId",
			fieldType: IntType,
		},
	}...)
	// City: id, name
	cities, err := mdl.NewTable("cities")
	if err != nil {
		return nil, err
	}

	cities.AddFields([]*Field{
		&Field{
			name:      "id",
			fieldType: IntType,
			pk:        true,
		},
		&Field{
			name:      "name",
			fieldType: StringType,
		},
	}...)
	return mdl, nil
}

// Helpers

func nPersonRecords(persons *Table, n int) ([]*InRecord, error) {
	records := make([]*InRecord, n)

	for i := 0; i < n; i++ {
		rec := persons.NewRecord()
		err := simplePersonRecord(rec, i)
		if err != nil {
			return nil, err
		}
		records[i] = rec
	}
	return records, nil
}

func simplePersonRecord(rec *InRecord, n int) error {
	rec.SetValue(FId, n)

	err := rec.SetValue(FName, "Fred_"+strconv.Itoa(n))
	if err != nil {
		return err
	}

	err = rec.SetValue(FAge, 54)
	if err != nil {
		return err
	}
	return nil
}

func two2PersonRecords2(sess *Session, persons *Table) error {
	err := sess.SaveFields(persons, []*Value{
		{field: persons.Field(FId), value: VPersonID0},
		{field: persons.Field(FName), value: VPersonName0},
		{field: persons.Field(FAge), value: VPersonAge0},
	})
	if err != nil {
		return err
	}
	err = sess.SaveFields(persons, []*Value{
		{field: persons.Field(FId), value: VPersonID1},
		{field: persons.Field(FName), value: VPersonName1},
		{field: persons.Field(FAge), value: VPersonAge1},
	})
	if err != nil {
		return err
	}

	return nil
}

func twoPersonRecords(persons *Table) ([]*InRecord, error) {
	records := make([]*InRecord, 2)

	rec1 := persons.NewRecord()
	records[0] = rec1
	err := rec1.SetValue(FId, VPersonID0)
	if err != nil {
		return nil, err
	}
	err = rec1.SetValue(FName, "Fred")
	if err != nil {
		return nil, err
	}

	err = rec1.SetValue(FAge, 54)
	if err != nil {
		return nil, err
	}

	rec2 := persons.NewRecord()
	records[1] = rec2
	err = rec2.SetValue(FId, VPersonID1)
	if err != nil {
		return nil, err
	}
	err = rec2.SetValue(FName, "Harry")
	if err != nil {
		return nil, err
	}

	err = rec2.SetValue(FAge, 21)
	if err != nil {
		return nil, err
	}

	return records, nil
}
