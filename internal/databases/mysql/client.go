package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	services = map[string]*gorm.DB{}
	name     = "mysql"
)

func Connect() {
	db, err := gorm.Open("mysql", "root:root@tcp(mysql:3306)/gin?charset=utf8")

	if err != nil {
		panic("failed to connect database")
	}

	//defer func() {
	//	err := db.Close()
	//	if err != nil {
	//		panic(err)
	//	}
	//}()

	//连接测试
	if err := db.DB().Ping(); err != nil {
		panic(err)
	}

	db.DB().SetMaxIdleConns(10)  //空闲连接
	db.DB().SetMaxOpenConns(100) //最大打开连接

	services[name] = db.New()
}

func Close() {
	if _, ok := services[name]; ok == true {
		err := services[name].Close()
		if err != nil {
			panic(err)
		}
	}
}

func Client() *gorm.DB {
	if _, ok := services[name]; ok == false {
		Connect()
	}

	return services[name]
}
