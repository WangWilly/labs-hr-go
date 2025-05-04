package testutils

import (
	"fmt"
	"testing"

	"gorm.io/gorm"
)

////////////////////////////////////////////////////////////////////////////////

func tableName(tx *gorm.DB, model interface{}) string {
	stmt := &gorm.Statement{DB: tx}
	stmt.Parse(model)
	return stmt.Schema.Table
}

func MustClearTable(t *testing.T, tx *gorm.DB, model interface{}) {
	tableName := tableName(tx, model)
	_, err := tx.Raw(fmt.Sprintf("DELETE FROM %s", tableName)).Rows()
	if err != nil {
		t.Fatal(err)
	}
}
