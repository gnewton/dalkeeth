package dalkeeth

import (
	"fmt"
	"log"
	"strconv"
	//	"os"
)

const TPerson = "persons"

// const TPersonK = "person_key"
const FId = "id"
const FName = "name"
const FNameDefaultValue = "no-name"
const FAge = "age"
const FAgeDefaultValue = "99"
const FAgeMinValue = 0
const FAgeMaxValue = 150
const FWeight = "weight"
const FWeightDefaultValue = "1"
const FCitizen = "citizen"
const FCitizenDefaultValue = "true"
const VPersonID0 = int64(43)
const VPersonName0 = "Fred"
const VPersonAge0 = 42
const VPersonWeight0 = 72

const VPersonID1 = int64(1090)
const VPersonName1 = "Sally"
const VPersonAge1 = 37
const VPersonWeight1 = 60

const TAddress = "addresses"

// const TAddressK = "address_key"
const FStreet = "street"
const FCity = "city"

const JTPersonName = "person_address"

// const JTPersonNameK = "person_address_key"
const FPersonId = "person_id"
const FAddressId = "address_id"

func testModel0() (*Model, error) {
	model := NewModel()

	persons, err := model.NewTable(TPerson)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	err = persons.AddFields([]*Field{
		&Field{
			name:      FId,
			fieldType: IntType,
			pk:        true,
		},
		&Field{
			name:         FAge,
			fieldType:    IntType,
			defaultValue: FAgeDefaultValue,
			rangge: &Range{
				min: FAgeMinValue,
				max: FAgeMaxValue,
			},
		},
		&Field{
			name:         FWeight,
			fieldType:    FloatType,
			defaultValue: FWeightDefaultValue,
		},
		&Field{
			name:         FCitizen,
			fieldType:    BoolType,
			defaultValue: FCitizenDefaultValue,
		},
		&Field{
			name:         FName,
			fieldType:    StringType,
			defaultValue: FNameDefaultValue,
		},
	}...)

	//
	addresses, err := model.NewTable(TAddress)
	if err != nil {
		return nil, err
	}
	if err = addresses.AddFields([]*Field{
		&Field{
			name:      FId,
			fieldType: IntType,
			pk:        true,
		},
		&Field{
			name:      FStreet,
			fieldType: StringType,
			length:    64,
			notNull:   true,
		},
		&Field{
			name:      FCity,
			fieldType: StringType,
			indexed:   true,
			notNull:   true,
			length:    64,
		}}...); err != nil {
		return nil, err
	}

	//
	person_address, err := model.NewTable(JTPersonName)
	if err != nil {
		return nil, err
	}
	if err = person_address.AddFields([]*Field{
		&Field{
			name:      FId,
			fieldType: IntType,
			pk:        true,
		},
		&Field{
			name:      FPersonId,
			fieldType: IntType,
			notNull:   true,
		},
		&Field{
			name:      FAddressId,
			fieldType: IntType,
			notNull:   true,
		}}...); err != nil {
		return nil, err
	}

	err = model.AddForeignKey(person_address, FPersonId, persons, FId)
	if err != nil {
		return nil, err
	}
	err = person_address.AddIndex(true, FPersonId, FAddressId)
	if err != nil {
		return nil, err
	}
	return model, model.Freeze()
}

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
		//fmt.Fprintln(os.Stdout, createSql)
		result, err := db.Exec(createSql)

		if err != nil {
			log.Println(fmt.Errorf("writeTestModelSchema: %s", err))
			return nil, err
		}
		_, err = result.RowsAffected()
		if err != nil {
			log.Println(fmt.Errorf("writeTestModelSchema: %s", err))
			return nil, err
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
