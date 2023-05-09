package dalkeeth

import (
	"fmt"
	"log"
	//	"os"
)

const TPerson = "persons"

// const TPersonK = "person_key"
const FId = "id"
const FName = "name"
const FAge = "age"
const FWeight = "weight"
const FCitizen = "citizen"
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

func defineTestModel() (*Model, error) {
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
			defaultValue: "42",
		},
		&Field{
			name:         FWeight,
			fieldType:    FloatType,
			defaultValue: "72",
		},
		&Field{
			name:         FCitizen,
			fieldType:    BoolType,
			defaultValue: "true",
		},
		&Field{
			name:         FName,
			fieldType:    StringType,
			defaultValue: "person's \"`name",
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

func initAndWriteTestTableSchema() (*Session, error) {
	model, err := defineTestModel()

	if err != nil {
		return nil, err
	}

	sess, err := NewSession(model)
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
			log.Println(fmt.Errorf("initAndWriteTestTableSchema: %s", err))
			return nil, err
		}
		_, err = result.RowsAffected()
		if err != nil {
			log.Println(fmt.Errorf("initAndWriteTestTableSchema: %s", err))
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

func fullModel() (*Model, error) {
	mdl := new(Model)
	//clients := mdl.NewTable("clients")
	clients, err := mdl.NewTable("clients")
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
