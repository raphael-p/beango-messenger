package database

import (
	"database/sql"
	"reflect"
)

type databaseEntity interface {
	User | Chat | ChatUser | MessageDatabase | Message
}

// Maps a SQL row onto a struct of a database entity
func scanRow[T databaseEntity](row *sql.Row) (*T, error) {
	target, scanArgs := prepForScan[T]()
	err := row.Scan(scanArgs...)
	return target, err
}

// Maps SQL rows onto a slice of structs of a database entity
func scanRows[T databaseEntity](rows *sql.Rows, err error) ([]T, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := []T{}
	for rows.Next() {
		target, scanArgs := prepForScan[T]()
		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}
		results = append(results, *target)
	}

	return results, rows.Err()
}

// Creates a pointer to a database entity and generates a slice of scan arguments from it
func prepForScan[T databaseEntity]() (*T, []any) {
	target := new(T)
	value := reflect.ValueOf(target).Elem()
	length := value.NumField()
	scanArgs := make([]any, length)
	for i := 0; i < length; i++ {
		scanArgs[i] = value.Field(i).Addr().Interface()
	}
	return target, scanArgs
}
