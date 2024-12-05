package db

import (
	"api-gateway/configs"
	"api-gateway/pkg/logger"
	"fmt"
	"time"

	"gorm.io/gorm"

	"gorm.io/driver/mysql"
)

type MySql struct {
	Client *gorm.DB
}

func NewMySqlClient(connString string, dbName string) *MySql {
	mysqlClient := ConnectDB(connString, dbName)
	return &MySql{Client: mysqlClient}
}

func ConnectDB(connString string, dbName string) *gorm.DB {

	checkDatabaseExists(connString, dbName)

	dsn := fmt.Sprintf("%s%s?charset=utf8mb4&parseTime=True&loc=Local",
		configs.C.MySql.ConnString,
		configs.C.MySql.Name,
	)

	// 連接數據庫

	clinet, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {

		logger.Fatal("連接數據庫失敗:", err)
	}

	sqlDB, err := clinet.DB()
	if err != nil {
		logger.Fatal("failed to get db instance: ", err)
	}
	sqlDB.SetMaxIdleConns(configs.C.MySql.MySqlBase.MaxIdleConns)
	sqlDB.SetMaxOpenConns(configs.C.MySql.MySqlBase.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(configs.C.MySql.MySqlBase.ConnMaxLifeTime * time.Minute)

	return clinet
}

func (m *MySql) CloseDB() {
	sqlDB, err := m.Client.DB()
	if err != nil {
		logger.Fatal("models.CloseDB err: ", err)
	}
	defer sqlDB.Close()
}

// CheckDatabaseExists checks if a database with the given name exists
func checkDatabaseExists(connString string, dbName string) {
	clinet, err := gorm.Open(mysql.Open(connString), &gorm.Config{})
	if err != nil {
		logger.Fatal("failed to connect database: %v", err)
	}

	var count int64
	clinet.Raw("SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = ?", dbName).Scan(&count)

	if count < 0 {
		result := clinet.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName))
		if result.Error != nil {
			logger.Fatal("failed to create database: ", result.Error)
		}
		logger.Info("Database %s created successfully\n", dbName)
	}

	sqlDB, err := clinet.DB()
	if err != nil {
		logger.Fatal("models.CloseDB err: ", err)
	}
	defer sqlDB.Close()
}

// func Migrate(dst ...interface{}) error {
// 	return clinet.AutoMigrate(dst...)
// }
