package db

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	xtime "github.com/xiaxin/moii/time"
	"time"
)

// TODO charset utf8mb4
const DsnStr = "%v:%v@tcp(%v:%v)/%v?parseTime=True&loc=Local&charset=utf8mb4"

type MysqlLogger interface {
	Print(v ...interface{})
}

type MysqlConfig struct {
	User     string `yaml:"user"`
	Host     string `yaml:"host"`
	Password string `yaml:"password"`
	Port     string `yaml:"port"`
	DbName   string `yaml:"dbname"`
	LogMode  bool   `yaml:"logmode"`

	Active      int            // pool
	Idle        int            // pool
	IdleTimeout xtime.Duration // connect max life time.
	Logger      MysqlLogger
}

func init() {
}

// NewMySQL new db and retry connection when has error.
func NewMysql(c *MysqlConfig) (*gorm.DB, error) {
	if nil == c {
		return nil, errors.New("mysql config is nil")
	}

	dsn := fmt.Sprintf(DsnStr, c.User, c.Password, c.Host, c.Port, c.DbName)

	db, err := gorm.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	//  设置日志级别
	db.LogMode(c.LogMode)

	//  设置日志
	if nil != c.Logger {
		db.SetLogger(c.Logger)
	}

	db.Set("gorm:table_options", "CHARSET=utf8mb4")

	//  TODO 这3个设置的详细说明
	db.DB().SetMaxIdleConns(c.Idle)

	db.DB().SetMaxOpenConns(c.Active)
	db.DB().SetConnMaxLifetime(time.Duration(c.IdleTimeout) * time.Second)

	return db, nil
}
