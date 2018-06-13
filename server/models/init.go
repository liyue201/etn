package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/liyue201/go-logger"
	"time"
)

var db *gorm.DB

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

func InitDb(user, pwd, host, port, dbname string) error {
	var err error

	dialect := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pwd, host, port, dbname)
	db, err = gorm.Open("mysql", dialect)
	if err != nil {
		logger.Errorf("Open db error [%s] with dialect[%s]", err.Error(), dialect)
		return err
	}

	db.DB().SetMaxOpenConns(50)
	db.DB().SetMaxIdleConns(5)
	db.DB().SetConnMaxLifetime(time.Second * 10)

	db.DB().Ping()
	if db.Error != nil {
		logger.Errorf("Ping db fail: %#v", db.Error)
		return db.Error
	}

	migrate()

	logger.Info("db inited")
	return nil
}

func CloseDB() {
	if db != nil {
		db.Close()
		db = nil
	}
}

func migrate()  {

	//创建表
	reslut := db.Set("gorm:table_options", "ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci").
		AutoMigrate(
		&File{},
	)
	if reslut.Error != nil {
		fmt.Printf("[migrate] %s\n", reslut.Error)
		return
	}
}
