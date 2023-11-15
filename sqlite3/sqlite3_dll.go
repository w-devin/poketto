package sqlite3

import "C"
import (
	"fmt"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

// SQLiteMsg sqlite3 error message: https://blog.csdn.net/czcdms/article/details/44461495
type SQLiteMsg int

const (
	SQLiteOK SQLiteMsg = iota
	SQLiteError
	SQLiteInternal
	SQLitePerm
	SQLiteAbort
	SQLiteBusy
	SQLiteLocked
	SQLiteNomem
	SQLiteReadonly
	SQLiteInterrupt
	SQLiteIOErr
	SQLiteCorrupt
	SQLiteNotfound
	SQLiteFull
	SQLiteCantopen
	SQLiteProtocol
	SQLiteEmpty
	SQLiteSchema
	SQLiteToobig
	SQLiteConstraint
	SQLiteMismatch
	SQLiteMisuse
	SQLiteNolfs
	SQLiteAuth
	SQLiteFormat
	SQLiteRange
	SQLiteNotadb
	SQLiteRow  = 100
	SQLiteDone = 101
)

const (
	SQLiteOKMsg         = "Successful result"
	SQLiteErrorMsg      = "SQL error or missing database"
	SQLiteInternalMsg   = "Internal logic error in SQLite"
	SQLitePermMsg       = "Access permission denied"
	SQLiteAbortMsg      = "Callback routine requested an abort"
	SQLiteBusyMsg       = "The database file is locked"
	SQLiteLockedMsg     = "A table in the database is locked"
	SQLiteNomemMsg      = "A malloc() failed"
	SQLiteReadonlyMsg   = "Attempt to write a readonly database"
	SQLiteInterruptMsg  = "Operation terminated by sqlite3_interrupt()"
	SQLiteIOErrMsg      = "Some kind of disk I/O error occurred"
	SQLiteCorruptMsg    = "The database disk image is malformed"
	SQLiteNotfoundMsg   = "Table or record not found"
	SQLiteFullMsg       = "Insertion failed because database is full"
	SQLiteCantopenMsg   = "Unable to open the database file"
	SQLiteProtocolMsg   = "Database lock protocol error"
	SQLiteEmptyMsg      = "Database is empty"
	SQLiteSchemaMsg     = "The database schema changed"
	SQLiteToobigMsg     = "String or BLOB exceeds size limit"
	SQLiteConstraintMsg = "Abort due to constraint violation"
	SQLiteMismatchMsg   = "Data type mismatch"
	SQLiteMisuseMsg     = "Library used incorrectly"
	SQLiteNolfsMsg      = "Uses OS features not supported on host"
	SQLiteAuthMsg       = "Authorization denied"
	SQLiteFormatMsg     = "Auxiliary database format error"
	SQLiteRangeMsg      = "2nd parameter to sqlite3_bind out of range"
	SQLiteNotadbMsg     = "File opened that is not a database file"
	SQLiteRowMsg        = "sqlite3_step() has another row ready"
	SQLiteDoneMsg       = "sqlite3_step() has finished executing"
	SQLiteUndefined     = "undefined"
)

const (
	SQLITE_TRANSIENT = 18446744073709551615

	SQLiteDataTypesInt   = 1
	SQLiteDataTypesFloat = 2
	SQLiteDataTypesText  = 3
	SQLiteDataTypesBlob  = 4
	SQLiteDataTypesNull  = 5
)

func (s SQLiteMsg) ErrCodeToMsg() string {
	switch s {
	case SQLiteOK:
		return SQLiteOKMsg
	case SQLiteError:
		return SQLiteErrorMsg
	case SQLiteInternal:
		return SQLiteInternalMsg
	case SQLitePerm:
		return SQLitePermMsg
	case SQLiteAbort:
		return SQLiteAbortMsg
	case SQLiteBusy:
		return SQLiteBusyMsg
	case SQLiteLocked:
		return SQLiteLockedMsg
	case SQLiteNomem:
		return SQLiteNomemMsg
	case SQLiteReadonly:
		return SQLiteReadonlyMsg
	case SQLiteInterrupt:
		return SQLiteInterruptMsg
	case SQLiteIOErr:
		return SQLiteIOErrMsg
	case SQLiteCorrupt:
		return SQLiteCorruptMsg
	case SQLiteNotfound:
		return SQLiteNotfoundMsg
	case SQLiteFull:
		return SQLiteFullMsg
	case SQLiteCantopen:
		return SQLiteCantopenMsg
	case SQLiteProtocol:
		return SQLiteProtocolMsg
	case SQLiteEmpty:
		return SQLiteEmptyMsg
	case SQLiteSchema:
		return SQLiteSchemaMsg
	case SQLiteToobig:
		return SQLiteToobigMsg
	case SQLiteConstraint:
		return SQLiteConstraintMsg
	case SQLiteMismatch:
		return SQLiteMismatchMsg
	case SQLiteMisuse:
		return SQLiteMisuseMsg
	case SQLiteNolfs:
		return SQLiteNolfsMsg
	case SQLiteAuth:
		return SQLiteAuthMsg
	case SQLiteFormat:
		return SQLiteFormatMsg
	case SQLiteRange:
		return SQLiteRangeMsg
	case SQLiteNotadb:
		return SQLiteNotadbMsg
	case SQLiteRow:
		return SQLiteRowMsg
	case SQLiteDone:
		return SQLiteDoneMsg
	default:
		return SQLiteUndefined
	}
}

var (
	sqlite3 = windows.NewLazyDLL(SQLITE3DLL)

	procSQLite3Open         = sqlite3.NewProc("sqlite3_open")
	procSQLite3Key          = sqlite3.NewProc("sqlite3_key")
	procSQLite3PrepareV2    = sqlite3.NewProc("sqlite3_prepare_v2")
	procSQLite3Step         = sqlite3.NewProc("sqlite3_step")
	procSQLite3ColumnCount  = sqlite3.NewProc("sqlite3_column_count")
	procSQLite3ColumnName   = sqlite3.NewProc("sqlite3_column_name")
	procSQLite3ColumnType   = sqlite3.NewProc("sqlite3_column_type")
	procSQLite3ColumnInt    = sqlite3.NewProc("sqlite3_column_int")
	procSQLite3ColumnDouble = sqlite3.NewProc("sqlite3_column_double")
	procSQLite3ColumnText   = sqlite3.NewProc("sqlite3_column_text")
	procSQLite3ColumnBlob   = sqlite3.NewProc("sqlite3_column_blob")
	procSQLite3Finalize     = sqlite3.NewProc("sqlite3_finalize")
)

func sqlite3_open(baseName string, database *uintptr) error {
	// reference: https://github.com/iamacarpet/go-sqlite3-win64/blob/master/sqlite3_raw.go#L71
	r1, _, _ := syscall.SyscallN(procSQLite3Open.Addr(), uintptr(unsafe.Pointer(C.CString(baseName))), uintptr(unsafe.Pointer(database)))
	if SQLiteMsg(r1) != SQLiteOK {
		return fmt.Errorf("failed to execute sqlite3_open, %s", SQLiteMsg(r1).ErrCodeToMsg())
	}

	return nil
}

func sqlite3_key(database *uintptr, key string) error {
	r1, _, _ := syscall.SyscallN(procSQLite3Key.Addr(), *database, uintptr(unsafe.Pointer(C.CString(key))), uintptr(len(key)))
	if SQLiteMsg(r1) != SQLiteOK {
		return fmt.Errorf("failed to execute sqlite3_key, %s", SQLiteMsg(r1).ErrCodeToMsg())
	}

	return nil
}

func sqlite3_prepare_v2(database *uintptr, query string) (uintptr, string, error) {
	var statement, excessData uintptr
	r1, _, _ := syscall.SyscallN(
		procSQLite3PrepareV2.Addr(),
		*database,
		uintptr(unsafe.Pointer(C.CString(query))),
		uintptr(len(query)),
		uintptr(unsafe.Pointer(&statement)),
		uintptr(unsafe.Pointer(&excessData)),
	)
	if SQLiteMsg(r1) != SQLiteOK {
		return 0, "", fmt.Errorf("failed to execute sqlite3_prepare_v2, %s", SQLiteMsg(r1).ErrCodeToMsg())
	}
	return statement, windows.BytePtrToString((*byte)(unsafe.Pointer(excessData))), nil
}

func sqlite3_step(statement uintptr) SQLiteMsg {
	r1, _, _ := syscall.SyscallN(
		procSQLite3Step.Addr(),
		statement,
	)
	return SQLiteMsg(r1)
}

func sqlite3_column_count(statement uintptr) (int, error) {
	r1, _, _ := syscall.SyscallN(
		procSQLite3ColumnCount.Addr(),
		statement,
	)

	return int(r1), nil
}

func sqlite3_column_name(statement uintptr, position int) string {
	r1, _, _ := syscall.SyscallN(
		procSQLite3ColumnName.Addr(),
		statement,
		uintptr(position),
	)

	return windows.BytePtrToString((*byte)(unsafe.Pointer(r1)))
}

func sqlite3_column_type(statement uintptr, position int) int {
	r1, _, _ := syscall.SyscallN(
		procSQLite3ColumnType.Addr(),
		statement,
		uintptr(position),
	)
	return int(r1)
}

func sqlite3_column_int(statement uintptr, position int) int {
	r1, _, _ := syscall.SyscallN(
		procSQLite3ColumnInt.Addr(),
		statement,
		uintptr(position),
	)
	return int(r1)
}

func sqlite3_column_double(statement uintptr, position int) float64 {
	r1, _, _ := syscall.SyscallN(
		procSQLite3ColumnDouble.Addr(),
		statement,
		uintptr(position),
	)
	return float64(r1)
}

func sqlite3_column_text(statement uintptr, position int) string {
	r1, _, _ := syscall.SyscallN(
		procSQLite3ColumnText.Addr(),
		statement,
		uintptr(position),
	)
	return windows.BytePtrToString((*byte)(unsafe.Pointer(r1)))
}

func sqlite3_column_blob(statement uintptr, position int) string {
	r1, _, _ := syscall.SyscallN(
		procSQLite3ColumnBlob.Addr(),
		statement,
		uintptr(position),
	)
	return windows.BytePtrToString((*byte)(unsafe.Pointer(r1)))
}

func sqlite3_exec() {
	// todo: need by ExecuteNonQuery
}

func sqlite3_finalize(statement uintptr) error {
	r1, _, _ := syscall.SyscallN(
		procSQLite3Finalize.Addr(),
		statement,
	)
	if SQLiteMsg(r1) != SQLiteOK {
		return fmt.Errorf("failed to execute sqlite3_finalize, %s", SQLiteMsg(r1).ErrCodeToMsg())
	}

	return nil
}
