package config

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db *gorm.DB

	username = "root"
	password = "example"
	host     = "localhost"
	port     = "3306"
	dbName   = "bookstore"
)

func Connect() {

	databaseConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbName)

	d, err := gorm.Open("mysql", databaseConnectionString)
	if err != nil {
		panic(err)
	}
	db = d
}

func GetDB() *gorm.DB {
	return db
}
