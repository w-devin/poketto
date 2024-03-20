package sqlite3

import "fmt"

type SQLiteBase struct {
	database uintptr
}

func OpenDatabase(baseName, dbKey string) (*SQLiteBase, error) {
	return nil, fmt.Errorf("only support windows")
}
