package db

import (
	"api-gateway/configs"
	"api-gateway/pkg/logger"
	"fmt"
	"time"

	"gorm.io/gorm"

	"gorm.io/driver/mysql"
)

var MySql *gorm.DB

func ConnectDB() {

	var err error

	user := configs.C.MySql.User         // 用戶名
	password := configs.C.MySql.Password // 密碼
	host := configs.C.MySql.Host         // 主機地址
	name := configs.C.MySql.Name         // 名稱

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/",
		user,
		password,
		host,
	)

	MySql, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect database: %v", err)
	}

	// Check if database exists, if not create it
	if !checkDatabaseExists(name) {
		createDatabase(name)
	}

	CloseDB()

	dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		host,
		name,
	)

	// 連接數據庫

	MySql, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {

		logger.Fatal("連接數據庫失敗:", err)
	}

	sqlDB, err := MySql.DB()
	if err != nil {
		logger.Fatal("failed to get db instance: ", err)
	}
	sqlDB.SetMaxIdleConns(configs.C.MySql.MySqlBase.MaxIdleConns)
	sqlDB.SetMaxOpenConns(configs.C.MySql.MySqlBase.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(configs.C.MySql.MySqlBase.ConnMaxLifeTime * time.Minute)

}

func CloseDB() {
	sqlDB, err := MySql.DB()
	if err != nil {
		logger.Fatal("models.CloseDB err: ", err)
	}
	defer sqlDB.Close()
}

// CheckDatabaseExists checks if a database with the given name exists
func checkDatabaseExists(dbName string) bool {
	var count int64
	MySql.Raw("SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", dbName).Scan(&count)
	return count > 0
}

// CreateDatabase creates a new database with the given name
func createDatabase(dbName string) {
	result := MySql.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
	if result.Error != nil {
		logger.Fatal("failed to create database: ", result.Error)
	}
	logger.Info("Database %s created successfully\n", dbName)
}

func Migrate(dst ...interface{}) error {
	return MySql.AutoMigrate(dst...)
}
