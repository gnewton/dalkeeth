package dalkeeth

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

// Assumed to work across all DBs ??

func recordExists(db *sql.DB, tableName string, id int64) (bool, error) {
	if db == nil {
		return false, errors.New("DB is nil")
	}
	if len(tableName) == 0 {
		return false, errors.New("Table name is empty string")
	}

	if id < 0 {
		return false, fmt.Errorf("Primary key id is < 0: %d", id)
	}

	q := "SELECT id from " + tableName + " where id=?"
	var value int64

	row := db.QueryRow(q, id)
	err := row.Scan(&value)
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return value == id, nil
}
