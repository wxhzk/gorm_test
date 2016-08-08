package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Profile   Profile
	ProfileID int
}

type Profile struct {
	gorm.Model
	Name string
}

func main() {
	db, err := gorm.Open("mysql", "root:mysql@/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err.Error())
	}
	//db.AutoMigrate(&User{}, &Profile{})
	db.Model(&User{}).Related(&Profile{})
}
