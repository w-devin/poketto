//go:build !windows
// +build !windows

package sqlite3

import (
	"fmt"
)

func OpenDatabase(baseName, dbKey string) (*SQLiteBase, error) {
	return nil, fmt.Errorf("only support windows")
}

type SQLiteBase struct {
	database uintptr
}

func (db *SQLiteBase) Close() {
}

func (db *SQLiteBase) ExecuteQuery(query string) (fields []string, ret []map[string]interface{}, err error) {
	return
}

func (db *SQLiteBase) GetResultFields(statement uintptr) ([]string, error) {
	return nil, nil
}

func (db *SQLiteBase) ReadNextRow(statement uintptr, fields *[]string) (map[string]interface{}, error) {
	return nil, nil
}
