package models

import (
	"fmt"
	mysql "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/liyue201/go-logger"
	"time"
	"usermgr/common/config"
)

var db *gorm.DB

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

func InitDb() error {
	var err error

	dialect := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Cfg.Db.User, config.Cfg.Db.Password, config.Cfg.Db.Host, config.Cfg.Db.Port, config.Cfg.Db.Dbname)
	db, err = gorm.Open(config.Cfg.Db.Driver, dialect)
	if err != nil {
		logger.Errorf("Open db error [%s] with dialect[%s]", err.Error(), dialect)
		return err
	}

	db.DB().SetMaxOpenConns(50)
	db.DB().SetMaxIdleConns(5)
	db.DB().SetConnMaxLifetime(time.Second * 10)
	db.SetLogger(DbLogger{})
	mysql.SetLogger(DbLogger{})

	db.DB().Ping()
	if db.Error != nil {
		logger.Errorf("Ping db fail: %#v", db.Error)
		return db.Error
	}

	//表名非复数形式
	//db.SingularTable(true)

	logger.Info("db inited")
	return nil
}

func CloseDB() {
	if db != nil {
		db.Close()
		db = nil
	}
}

// logger
type DbLogger struct {
	gorm.LogWriter
}

// Print format & print log
func (l DbLogger) Print(values ...interface{}) {
	logger.Info(values)
}

func TransactionBegin() *gorm.DB {
	return db.Begin()
}

func TransactionCommit(tx *gorm.DB) {
	tx.Commit()
}

func InitTestDb(dbname, pwd, host, port string) {
	var err error
	dialect := fmt.Sprintf("root:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", pwd, host, port, dbname)
	db, err = gorm.Open("mysql", dialect)
	if err != nil {
		fmt.Printf("Open db error [%s] with dialect[%s]", err.Error(), dialect)
		panic(err)
	}

	db, err = gorm.Open("mysql", dialect)
	if err != nil {
		fmt.Printf("Open db error [%s] with dialect[%s]", err.Error(), dialect)
		panic(err)
	}
}
