package dalkeeth

import ()

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
const FPersonId = "person_id"
const FAddressId = "address_id"

var XPersonIdField *Field
var XAddressField *Field

func testModel0() (*Model, error) {
	model := NewModel()

	persons, err := model.NewTable(TPerson)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	XPersonIdField = &Field{
		name:      FId,
		fieldType: IntType,
		pk:        true,
	}

	err = persons.AddFields([]*Field{
		XPersonIdField,
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
