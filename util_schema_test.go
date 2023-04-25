package dalkeeth

import (
	"fmt"
	"log"
)

const TPerson = "persons"
const TPersonK = "person_key"
const FId = "id"
const FName = "name"
const FAge = "age"
const FWeight = "weight"
const FCitizen = "citizen"
const VPersonID0 = int64(43)
const VPersonID1 = int64(89)

const TAddress = "addresses"
const TAddressK = "address_key"
const FStreet = "street"
const FCity = "city"

const JTPersonName = "person_address"
const JTPersonNameKey = "person_address_key"
const FPersonId = "person_id"
const FAddressId = "address_id"

func defineTestModel() (*Model, error) {
	//mgr := NewSession()
	model := NewModel()

	persons, err := NewTable(TPerson)
	if err != nil {
		return nil, err
	}
	model.AddTable(TPerson, persons)

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

	err = model.AddTable(TPersonK, persons) // FIXXX
	if err != nil {
		return nil, err
	}

	//
	addresses, err := NewTable(TAddress)
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
	err = model.AddTable(TAddressK, addresses) // FIXXX
	if err != nil {
		return nil, err
	}
	//
	person_address, err := NewTable(JTPersonName)
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

	err = model.AddTable(JTPersonNameKey, person_address) // FIXXX
	if err != nil {
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

	return model, nil
}

func initAndPopulateTestTables() (*Session, error) {
	model, err := defineTestModel()

	if err != nil {
		return nil, err
	}

	mgr := NewSessionWithModel(model)
	mgr.dialect = new(DialectSqlite3)
	sqls, err := mgr.CreateTablesSQL()

	if err != nil {
		return nil, err
	}

	indexesSql, err := mgr.CreateTableIndexesSQL()

	if err != nil {
		return nil, err
	}

	db, err := openTestDB()
	if err != nil {
		return nil, err
	}

	mgr.db = db
	mgr.dialect = new(DialectSqlite3)

	// Create tables sql
	log.Println("Sql tables:", sqls)
	for i := 0; i < len(sqls); i++ {
		createSql := sqls[i]
		result, err := db.Exec(createSql)

		if err != nil {
			log.Println(fmt.Errorf("initAndPopulateTestTables: %s", err))
			return nil, err
		}
		_, err = result.RowsAffected()
		if err != nil {
			log.Println(fmt.Errorf("initAndPopulateTestTables: %s", err))
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

	return mgr, nil
}
