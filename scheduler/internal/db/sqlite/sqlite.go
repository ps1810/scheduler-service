package sqlite

import (
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"path/filepath"
	"scheduler/config"
	"scheduler/internal/logger"
	"sync"
	"time"
)

var DB *gorm.DB
var m sync.Mutex

func GetSqliteDB() *gorm.DB {
	if DB == nil {
		m.Lock()
		defer m.Unlock()

		logger.Log.Info("Initializing db again")
		err := InitSqliteDatabase(config.GetConfig().Sqlite)
		if err != nil {
			logger.Log.Error("Failed to initialize database", zap.Error(err))
		}
		logger.Log.Info("Database initialized")
	}
	return DB
}

func InitSqliteDatabase(sqliteConfig config.Sqlite) error {
	dbPath := filepath.Join(sqliteConfig.Path, sqliteConfig.Name)
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(sqliteConfig.MaxConnections)
	sqlDB.SetMaxIdleConns(sqliteConfig.MaxConnIdleTime)
	sqlDB.SetConnMaxLifetime(time.Hour)

	DB = db
	return nil
}
