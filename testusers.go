package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func init() {
	db_test, err := gorm.Open("mysql", "root:mysql@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True")
	if err != nil {
		panic(err.Error())
	}
	db = db_test
}

func main() {
	if db.HasTable("users") {
		fmt.Println("table users has already exists!")
	} else {
		fmt.Println("table users is not exists!")
	}
}
