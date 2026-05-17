package database

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func OpenReadOnlySQLite(path string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("file:%s?mode=ro&_busy_timeout=5000", path)
	return openSQLite(dsn)
}

// OpenReadWriteSQLite 打开读写模式的数据库连接
func OpenReadWriteSQLite(path string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("file:%s?mode=rw&_busy_timeout=5000", path)
	return openSQLite(dsn)
}

// OpenSQLiteCreate 打开数据库连接，如果不存在则创建
func OpenSQLiteCreate(path string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("file:%s?_busy_timeout=5000", path)
	return openSQLite(dsn)
}

func openSQLite(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(1)
	sqlDB.SetMaxOpenConns(1)

	return db, nil
}
