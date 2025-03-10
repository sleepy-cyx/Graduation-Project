package repo

import (
	"Graduation-Project/config"
	"Graduation-Project/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var DB *gorm.DB

func InitDB() (*gorm.DB, error) {
	dsn := config.Conf.MysqlDsn
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Logger.Errorf("[InitDB] err: %v", err)
		return nil, err
	}
	sqlDB, err := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return DB, nil
}
