package assets

import "embed"

var (
	//go:embed sqlite3.dll
	sqlite3Dll embed.FS
)

func GetSqliteDll() embed.FS {
	return sqlite3Dll
}
