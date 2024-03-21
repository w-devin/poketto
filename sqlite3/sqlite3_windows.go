package sqlite3

import (
	"C"
	"fmt"
	"github.com/w-devin/poketto/assets"
	"github.com/w-devin/poketto/file"
	"log"
	"os"
	"sync"
)

const (
	SQLITE3DLL = "sqlite3.dll"
)

var (
	instanceCountLock sync.Mutex
	instanceCount     = 0
)

type SQLiteBase struct {
	database uintptr
}

func Init() {
	if !file.IsFileExists(SQLITE3DLL) {
		sqlite3Dll := assets.GetSqliteDll()
		sqliteDll, err := sqlite3Dll.ReadFile(SQLITE3DLL)
		if err != nil {
			log.Fatalf("failed to got %s, %v", SQLITE3DLL, err)
		}

		err = os.WriteFile(SQLITE3DLL, sqliteDll, os.ModePerm)
		if err != nil {
			log.Fatalf("failed to put %s done, %v", SQLITE3DLL, err)
		}
	}
}

func OpenDatabase(baseName, dbKey string) (*SQLiteBase, error) {
	Init()
	instanceCountLock.Lock()
	instanceCount += 1
	instanceCountLock.Unlock()

	db := &SQLiteBase{}
	err := sqlite3_open(baseName, &db.database)
	if err != nil {
		return nil, fmt.Errorf("failed to execute sqlite3_open, %v", err)
	}

	if len(dbKey) != 0 {
		err := sqlite3_key(&db.database, dbKey)
		if err != nil {
			return nil, fmt.Errorf("failed to execute sqlite3_key, %v", err)
		}
	}

	return db, nil
}

func (db *SQLiteBase) Close() {
	instanceCountLock.Lock()
	instanceCount -= 1
	instanceCountLock.Unlock()
	if instanceCount <= 0 {
		if file.IsFileExists(SQLITE3DLL) {
			_ = os.Remove(SQLITE3DLL)
		}
	}
}

func (db *SQLiteBase) ExecuteQuery(query string) (fields []string, ret []map[string]interface{}, err error) {
	statement, _, err := sqlite3_prepare_v2(&db.database, query)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to execute sqlite3_prepare_v2, %v", err)
	}
	defer sqlite3_finalize(statement)

	for {
		row, err := db.ReadNextRow(statement, &fields)
		if err != nil {
			return fields, ret, fmt.Errorf("failed to Read New Row, %v", err)
		}

		if row == nil {
			break
		}
		ret = append(ret, row)
	}

	return fields, ret, nil
}

func (db *SQLiteBase) executeNonQuery(query string) {
	//  todo: use sqlte3_exec
	/*
		IntPtr error;
		sqlite3_exec(database, StringToPointer(query), IntPtr.Zero, IntPtr.Zero, out error);
		if (error != IntPtr.Zero)
			throw new Exception("Error with executing non-query: \"" + query + "\"!\n" + PointerToString(sqlite3_errmsg(error)));
	*/
	return
}

func (db *SQLiteBase) GetResultFields(statement uintptr) ([]string, error) {
	columnCount, err := sqlite3_column_count(statement)
	if err != nil {
		return nil, err
	}

	var columnNames []string
	for i := 0; i < columnCount; i++ {
		columnNames = append(columnNames, sqlite3_column_name(statement, i))
	}

	return columnNames, nil
}

func (db *SQLiteBase) ReadNextRow(statement uintptr, fields *[]string) (map[string]interface{}, error) {
	resultType := sqlite3_step(statement)
	if resultType != SQLiteRow {
		return nil, nil
	}

	ret := make(map[string]interface{})
	var err error
	if fields == nil || len(*fields) == 0 {
		*fields, err = db.GetResultFields(statement)
		if err != nil {
			return nil, err
		}
	}

	for i := 0; i < len(*fields); i++ {
		columnType := sqlite3_column_type(statement, i)
		switch columnType {
		case SQLiteDataTypesInt:
			ret[(*fields)[i]] = sqlite3_column_int(statement, i)
		case SQLiteDataTypesFloat:
			ret[(*fields)[i]] = sqlite3_column_double(statement, i)
		case SQLiteDataTypesText:
			ret[(*fields)[i]] = sqlite3_column_text(statement, i)
		case SQLiteDataTypesBlob:
			ret[(*fields)[i]] = sqlite3_column_blob(statement, i)
		default:
			ret[(*fields)[i]] = ""
		}
	}

	return ret, nil
}
