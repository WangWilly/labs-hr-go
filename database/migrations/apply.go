package migrations

import (
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/gorm"
)

var migrationList = []*gormigrate.Migration{
	{
		ID: "00001",
		Migrate: func(tx *gorm.DB) error {
			return Up00001Init(tx)
		},
		Rollback: func(tx *gorm.DB) error {
			return Down00001Init(tx)
		},
	},
}

func Apply(db *gorm.DB) error {
	m := gormigrate.New(db, gormigrate.DefaultOptions, migrationList)
	if err := m.Migrate(); err != nil {
		return err
	}

	return nil
}
