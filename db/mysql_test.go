package db

import (
	"github.com/xiaxin/moii/log"
	"testing"
)

type TestTable struct {
	ID int `gorm:"column:user_id; UNIQUE_INDEX"`
}

func TestMysql(t *testing.T) {
	log := log.Named("test")

	db := NewMysql(&MysqlConfig{
		User:        "root",
		Host:        "127.0.0.1",
		Password:    "880728",
		Port:        "3306",
		DbName:      "test",
		LogMode:     true,
		Active:      0,
		Idle:        10, // 空闲链接
		IdleTimeout: 5,
		Logger: log.NewGorm(),
	})

	if err := db.DropTableIfExists(&TestTable{}).Error; nil != err {
		log.Error(err)
	}

	if err := db.AutoMigrate(&TestTable{}).Error; nil != err {
		log.Error(err)
	}
}
