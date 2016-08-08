package main

import (
	"fmt"
	"time"

	//_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func init() {
	//此处设置的charset=utf8貌似没起效果，待查
	db_test, err := gorm.Open("mysql", "root:mysql@tcp(127.0.0.1:3306)/test?charset=utf8&parseTime=True")
	if err != nil {
		panic(err.Error())
	}
	db = db_test
}

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func (self *Product) String() string {
	return fmt.Sprintf("[%d-%s](%d)", self.ID, self.Code, self.Price)
}

func testProduct() {
	if db.HasTable(&Product{}) {
		fmt.Println("table products is already exists!")
	}
	db.AutoMigrate(&Product{})
	pro1 := &Product{Code: "Java", Price: 100}
	fmt.Println(pro1) //[0-Java](100)
	db.Create(pro1)
	fmt.Println(pro1) //[1-Java](100)

	pro2 := &Product{Code: "Python", Price: 300}
	fmt.Println(pro2)
	db.Save(pro2)
	fmt.Println(pro2)

	var pro3 Product
	db.First(&pro3, "code = ?", "Java")
	fmt.Println(pro3)
}

/*
type User struct {
	Id        uint64
	Name      string
	Password  string
	Salt      string
	CreatedAt uint64
	UpdatedAt uint64
}
*/

type User struct {
	Id              int
	Birthday        time.Time
	Age             int
	Name            string `gorm:"size:255"`
	CreditCard      CreditCard
	Emails          []Email
	BillingAddress  Address
	ShippingAddress Address
	Languages       []Language `gorm:"many2many:user_languages;"`
}

type Email struct {
	Id         int    `gorm:"primary_key;auto_increment"`
	UserId     int    `gorm:"index"`
	Email      string `gorm:"type:varchar(100);unique"`
	Subscribed bool
}

type Address struct {
	Id       int
	Address1 string `gorm:"not null;unique"`
	Address2 string `gorm:type:varchar(100);unique`
}

type Language struct {
	Id   int
	Name string `gorm:"index:index_name_code"`
	Code string `gorm:"index:index_name_code"`
}

type CreditCard struct {
	gorm.Model
	UserId int
	Number string
}

func testUser() {
	db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8")
	db.Debug()
	db.DropTableIfExists(&User{}, &Email{}, &Address{}, &Language{}, &CreditCard{})
	db.AutoMigrate(&Email{}, &Address{}, &Language{}, &CreditCard{}, &User{})
	//db.CreateTable(&User{})
	if db.HasTable(&User{}) {
		fmt.Println("table users has already exists!")
	} else {
		fmt.Println("create table users")
	}
}

func main() {
	testUser()
	//db.DropTable("users", "products")
}
