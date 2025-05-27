package cmd

import (
	"context"
	"go.uber.org/zap"
	"log"
	"scheduler/config"
	"scheduler/internal/app/scheduler"
	"scheduler/internal/db/sqlite"
	"scheduler/internal/local_cron"
	"scheduler/internal/logger"
	"scheduler/internal/repository"
	"time"
)

const defaultConfigFile = "config/config.yaml"

// SetUpConfig Reading the config file and loading it using viper
func SetUpConfig() {
	log.Default().Printf("Using config file: %s", defaultConfigFile)
	config.SetConfig(defaultConfigFile)
}

// SetUpLogger Loading the zap logger
func SetUpLogger() {
	log.Default().Printf("Using log level: %s", config.GetConfig().Log.Level)
	logger.InitLogger("zap")
}

// SetUpDatabase Intiailizing the database. Currently it is using sqlite
func SetUpDatabase() {
	logger.Log.Info("Initializing database connection")
	err := sqlite.InitSqliteDatabase(config.GetConfig().Sqlite)
	if err != nil {
		logger.Log.Fatal("Failed to initialize the database", zap.Error(err))
	}
	logger.Log.Info("Database initialized")
}

// SetUpCron Setting up the cron to run the scheduled jobs
func SetUpCron() {
	logger.Log.Info("Initializing Cron")
	local_cron.InitCron()
	local_cron.StartCron()
	logger.Log.Info("Cron Initialized")
}

// LoadSchedules Loading the existing schedules
func LoadSchedules() {
	logger.Log.Info("Loading schedules from database")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	repo := repository.NewRepository()
	schedulerApp := scheduler.NewAppScheduler(repo)
	err := schedulerApp.LoadAndScheduleJobs(ctx)
	if err != nil {
		logger.Log.Error("Unable to load existing schedules")
	}
}

func ShutDown() {
	logger.Log.Info("Shutting down cron")
	local_cron.StopCron()
}
