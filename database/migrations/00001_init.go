package migrations

import (
	"context"

	"github.com/WangWilly/labs-hr-go/pkgs/models"
	"github.com/WangWilly/labs-hr-go/pkgs/utils"
)

func Up00001Init(ctx context.Context) error {
	// This code is executed when the migration is applied.
	db, err := utils.GetCtxDb(ctx)
	if err != nil {
		return err
	}

	// Create the employeeinfo table
	if ok := db.Migrator().HasTable(&models.EmployeeInfo{}); !ok {
		if err := db.Migrator().CreateTable(&models.EmployeeInfo{}); err != nil {
			return err
		}
	}
	// Create the employeeattendance table
	if ok := db.Migrator().HasTable(&models.EmployeeAttendance{}); !ok {
		if err := db.Migrator().CreateTable(&models.EmployeeAttendance{}); err != nil {
			return err
		}
	}
	// Create the employeeinfo table
	if ok := db.Migrator().HasTable(&models.EmployeePosition{}); !ok {
		if err := db.Migrator().CreateTable(&models.EmployeePosition{}); err != nil {
			return err
		}
	}

	return nil
}

func Down00001Init(ctx context.Context) error {
	// This code is executed when the migration is rolled back.
	db, err := utils.GetCtxDb(ctx)
	if err != nil {
		return err
	}

	// Drop the employeeinfo table
	if ok := db.Migrator().HasTable(&models.EmployeeInfo{}); !ok {
		if err := db.Migrator().DropTable(&models.EmployeeInfo{}); err != nil {
			return err
		}
	}
	// Drop the employeeattendance table
	if ok := db.Migrator().HasTable(&models.EmployeeAttendance{}); !ok {
		if err := db.Migrator().DropTable(&models.EmployeeAttendance{}); err != nil {
			return err
		}
	}
	// Drop the employeeinfo table
	if ok := db.Migrator().HasTable(&models.EmployeePosition{}); !ok {
		if err := db.Migrator().DropTable(&models.EmployeePosition{}); err != nil {
			return err
		}
	}

	return nil
}
