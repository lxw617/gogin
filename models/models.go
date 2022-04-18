package models

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"time"

	"gogin/pkg/setting"
)

var DB *gorm.DB //全局数据库变量

type Model struct {
	//所有数据库表的公用字段
	ID        int       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func init() {
	//初始化变量
	var (
		err                                               error
		dbType, dbName, user, password, host, tablePrefix string
	)

	//获取配置文件中名为 database 的信息
	sec, err := setting.Cfg.GetSection("database")
	if err != nil {
		log.Fatal(2, "Fail to get section 'database': %v", err)
	}

	//获取指定配置信息
	dbType = sec.Key("TYPE").String()
	dbName = sec.Key("NAME").String()
	user = sec.Key("USER").String()
	password = sec.Key("PASSWORD").String()
	host = sec.Key("HOST").String()
	tablePrefix = sec.Key("TABLE_PREFIX").String() //blog_

	//连接数据库
	DB, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName))

	//如果错误打印错误
	if err != nil {
		log.Println(err)
	}

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}

	//打印执行sql
	DB.LogMode(true)
	DB.SingularTable(true)

	//GORM 使用 database/sql 维护连接池
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	DB.DB().SetMaxIdleConns(10)
	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	DB.DB().SetMaxOpenConns(100)
	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	//db.DB().SetConnMaxLifetime(time.Hour)
}

func CloseDB() {
	defer DB.Close()
}
