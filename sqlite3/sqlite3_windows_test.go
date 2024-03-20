package sqlite3

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOpenDatabase(t *testing.T) {
	key := "3f17fa99-9804-4189-a75f-39589413f94f"
	db, err := OpenDatabase("../test/assis2.db", key)
	assert.NoErrorf(t, err, "failed to open db, %v", err)
	defer db.Close()

	fields, datas, err := db.ExecuteQuery("select * from tb_account")
	assert.NoErrorf(t, err, "failed to execute query, %v", err)

	for _, data := range datas {
		for _, key := range fields {
			fmt.Printf("%s: %v, ", key, data[key])
		}
		fmt.Println()
	}
}
