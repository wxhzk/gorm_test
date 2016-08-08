package main

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	db *gorm.DB
)

func init() {
	db_test, err := gorm.Open("mysql", "root:mysql@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True&autocommit=true")
	if err != nil {
		panic(fmt.Sprintf("failed to connect database, error:%s", err.Error()))
	}
	db = db_test
}

//成员只支持基础类型、slice、struct；如果把下面的[]Email改成[]*Email，将提示警告信息
type User struct {
	Id         uint64
	Account    string     `gorm:"type:varchar(64);not null;default '';unique"`
	Passwd     string     `gorm:"type:varchar(64);not null;default '' "`
	Type       uint32     `gorm:"not null;default 0;"`
	Emails     []Email    //`gorm:"ForeignKey:UserId"` //不会主动加外键，需要自己添加
	Languages  []Language `gorm:"many2many:user_languages"`
	CreditCard CreditCard
}

type Email struct {
	Id     uint64
	UserId uint64 `gorm:"column:user_id;index"`
	Email  string `gorm:"type:varchar(128);not null;default '';unique"`
}

type Language struct {
	Id   uint64 `gorm:"primary_key;auto_increment"`
	Name string `gorm:"type:varchar(128);not null;default '';unique_index:index_name_code"`
	Code string `gorm:"type:varchar(128);not null;default '';unique_index:index_name_code"`
}

type CreditCard struct {
	Id     uint64
	UserId uint64 `gorm:"not null;default 0"`
	Number string `gorm:"type:varchar(128);not null;default '';unique;"`
}

func test() {
	db.DropTableIfExists(&Email{}, &User{}, &CreditCard{}, &Language{})
	//db.Set("gorm:table_options", "ENGINE=MyISAM CHARSET=utf8 COLLATE=utf8_unicode_ci")
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8 COLLATE=utf8_unicode_ci").AutoMigrate(&User{}, &CreditCard{}, &Email{}, &Language{})
	db.Model(&Email{}).ModifyColumn("user_id", "bigint(20) unsigned not null default 0")
	//外键可不添加
	//db.Model(&Email{}).AddForeignKey("user_id", "users(id)", "RESTRICT", "RESTRICT")
	if db.HasTable(&User{}) {
		fmt.Println("table users is create successful!")
	}
	u := &User{Account: "nnn", Passwd: "123456"}
	e := &Email{Id: 1001, UserId: 1, Email: "test email"}
	l := &Language{Name: "zh-cn", Code: "Chinese"}
	u.Emails = append(u.Emails, *e)
	u.Languages = append(u.Languages, *l)
	u.CreditCard.Number = "12345678"

	//这个时候u还没有Id，没有相关联的数据，所以此处添加无效
	//db.Model(u).Association("Languages").Append(Language{Name: "en", Code: "en"})
	fmt.Println("create")
	db.Create(u) //会自动保存emails,languages,user_languages
	db.Create(e) //如果主键有赋值，则需要手动保存数据库
	//此处会自动保存,Delete/Replace/Clear/Count
	db.Model(u).Association("Languages").Append(Language{Name: "en", Code: "en"})
	fmt.Printf("%+v\n", u)
	db.Create(&Email{UserId: 1, Email: "test1"}) //没有关联到宿主的也需要手动保存
	u1 := User{Id: 1}
	db.First(&u1, 1)

	c := db.Model(&u1).Association("Languages").Count()
	fmt.Println(c)
	db.Model(&u1).Related(&u1.CreditCard)
	db.Model(&u1).Related(&u1.CreditCard, "CreditCard")
	db.Model(&u1).Association("Languages").Find(&u1.Languages) //多对多的麻烦些
	//db.Model(&u1).Association("Emails").Find(&u1.Emails)
	//db.Model(&u1).Related(&u1.Emails) //功能同上
	db.Find(&u1.Emails, "user_id=?", u1.Id)
	fmt.Printf("%+v\n", u1)
	var l2 Language
	db.Where("code=?", "en").FirstOrInit(&l2) //{2,en,en}
	fmt.Println(l2)                           //{2,en,en}
	l2.Id = 0
	db.Where("code=?", "en").FirstOrInit(&l2, Language{Id: 5, Code: "jp", Name: "jp"})
	fmt.Println(l2) //{5,jp,jp}
}

func main() {
	test()
}
