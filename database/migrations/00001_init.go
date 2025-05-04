package migrations

import (
	"github.com/WangWilly/labs-hr-go/pkgs/models"
	"gorm.io/gorm"
)

func Up00001Init(db *gorm.DB) error {
	// This code is executed when the migration is applied.

	// Create the employeeinfo table
	if err := db.Migrator().CreateTable(&models.EmployeeInfo{}); err != nil {
		return err
	}
	// Create the employeeattendance table
	if err := db.Migrator().CreateTable(&models.EmployeeAttendance{}); err != nil {
		return err
	}
	// Create the employeeposition table
	if err := db.Migrator().CreateTable(&models.EmployeePosition{}); err != nil {
		return err
	}

	return nil
}

func Down00001Init(db *gorm.DB) error {
	// This code is executed when the migration is rolled back.

	// Drop the employeeinfo table
	if err := db.Migrator().DropTable(&models.EmployeeInfo{}); err != nil {
		return err
	}
	// Drop the employeeattendance table
	if err := db.Migrator().DropTable(&models.EmployeeAttendance{}); err != nil {
		return err
	}

	// Drop the employeeposition table
	if err := db.Migrator().DropTable(&models.EmployeePosition{}); err != nil {
		return err
	}

	return nil
}
