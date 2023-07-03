package database

import (
	"database/sql"
	"reflect"
)

type databaseEntity interface {
	User | Chat | ChatUser | Message
}

// Maps a SQL row onto a struct of a database entity
func scanRow[T databaseEntity](row *sql.Row) (*T, error) {
	target := new(T)
	value := reflect.ValueOf(target).Elem()
	length := value.NumField()
	scanArgs := make([]any, length)
	for i := 0; i < length; i++ {
		scanArgs[i] = value.Field(i).Addr().Interface()
	}
	err := row.Scan(scanArgs...)
	return target, err
}
