package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WangWilly/labs-hr-go/controllers/attendance"
	"github.com/WangWilly/labs-hr-go/controllers/employee"
	"github.com/WangWilly/labs-hr-go/database/migrations"
	"github.com/WangWilly/labs-hr-go/pkgs/repos/employeeattendancerepo"
	"github.com/WangWilly/labs-hr-go/pkgs/repos/employeeinforepo"
	"github.com/WangWilly/labs-hr-go/pkgs/repos/employeepositionrepo"
	"github.com/WangWilly/labs-hr-go/pkgs/seed"
	"github.com/WangWilly/labs-hr-go/pkgs/timemodule"
	"github.com/WangWilly/labs-hr-go/pkgs/utils"

	"github.com/sethvargo/go-envconfig"
)

////////////////////////////////////////////////////////////////////////////////

type envConfig struct {
	Port     string         `env:"PORT,default=8080"`
	Host     string         `env:"HOST,default=0.0.0.0"`
	DbCfg    utils.DbConfig `env:",prefix="`
	WithSeed bool           `env:"WITH_SEED,default=false"`
}

////////////////////////////////////////////////////////////////////////////////

func main() {
	ctx := context.Background()

	// Load environment variables
	cfg := &envConfig{}
	err := envconfig.Process(ctx, cfg)
	if err != nil {
		panic(err)
	}

	////////////////////////////////////////////////////////////////////////////
	// Initialize Gin router

	r := utils.GetDefaultRouter()

	////////////////////////////////////////////////////////////////////////////
	// Setup database

	db, err := utils.GetDB(cfg.DbCfg)
	if err != nil {
		panic(err)
	}

	if err := migrations.Apply(ctx, db); err != nil {
		panic(err)
	}

	// Seed the database if WITH_SEED is true
	if cfg.WithSeed {
		fmt.Println("Seeding database with dummy data...")
		if err := seed.SeedData(ctx, db); err != nil {
			panic(fmt.Errorf("failed to seed database: %w", err))
		}
		fmt.Println("Database seeded successfully!")
	}

	////////////////////////////////////////////////////////////////////////////

	timeModule := timemodule.New()
	employeeInfoRepo := employeeinforepo.New()
	employeePositionRepo := employeepositionrepo.New()
	employeeAttendanceRepo := employeeattendancerepo.New()

	////////////////////////////////////////////////////////////////////////////
	// Initialize the controllers

	employeeCtrlCfg := employee.Config{}
	employeeCtrl := employee.NewController(employeeCtrlCfg, db, timeModule, employeeInfoRepo, employeePositionRepo)
	employeeCtrl.RegisterRoutes(r)

	attendanceCtrlCfg := attendance.Config{}
	attendanceCtrl := attendance.NewController(attendanceCtrlCfg, db, timeModule, employeePositionRepo, employeeAttendanceRepo)
	attendanceCtrl.RegisterRoutes(r)

	////////////////////////////////////////////////////////////////////////////

	// Set up the server
	srv := &http.Server{
		Addr:    cfg.Host + ":" + cfg.Port,
		Handler: r,
	}

	// Start the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	////////////////////////////////////////////////////////////////////////////

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	// Kill (no param) default sends syscall.SIGTERM
	// Kill -2 is syscall.SIGINT
	// Kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Log shutdown message
	fmt.Println("Received shutdown signal, shutting down server...")
	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// TODO:

	// Gracefully shutdown the server
	fmt.Println("Shutdown HTTP Server ...")
	if err := srv.Shutdown(ctx); err != nil {
		// Handle shutdown error
		panic(err)
	}

	// Wait for tasks to finish or timeout
	<-ctx.Done()
	fmt.Println("Server shutdown complete.")
}
