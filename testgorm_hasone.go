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

type User struct {
	Id         uint64
	Name       string     `gorm:"type:varchar(128);unique;default:''"`
	CreditCard CreditCard `gorm:"ForeignKey:UserId;AssociationForeignKey:Refer"`
}

type CreditCard struct {
	gorm.Model
	UserId uint64
	Number string `gorm:"type:varchar(255);unique;default:''"`
}

func test() {
	db.DropTableIfExists(&User{}, &CreditCard{})
	db.AutoMigrate(&User{}, &CreditCard{})
	u := &User{1001, "nnn", CreditCard{Number: "123456"}}
	err := db.Save(u).Error
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(u)
}

func main() {
	test()
}
