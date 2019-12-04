package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"strconv"
	"time"
)

type Person struct {
	Id        int    `gorm:"column:id; PRIMARY_KEY ;AUTO_INCREMENT" json:"id" form:"id"`
	FirstName string `gorm:"column:first_name" json:"first_name" form:"first_name"`
	LastName  string `gorm:"column:last_name" json:"last_name" form:"last_name"`
}

func main() {
	db, err := gorm.Open("mysql", "root:root@tcp(mysql:3306)/gin?charset=utf8")

	if err != nil {
		panic("failed to connect database")
	}

	defer func() {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}()

	//连接测试
	if err := db.DB().Ping(); err != nil {
		panic(err)
	}

	db.DB().SetMaxIdleConns(10)  //空闲连接
	db.DB().SetMaxOpenConns(100) //最大打开连接

	gin.SetMode(gin.DebugMode)

	router := gin.Default()

	// 首页
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "It works")
	})

	// 增
	router.POST("/person", func(c *gin.Context) {
		firstName := c.Request.FormValue("first_name")
		lastName := c.Request.FormValue("last_name")

		rs, err := db.DB().Exec("INSERT INTO person(first_name, last_name) VALUES (?, ?)", firstName, lastName)
		if err != nil {
			panic(err)
		}

		id, err := rs.LastInsertId()
		if err != nil {
			panic(err)
		}

		msg := fmt.Sprintf("insert successful %d", id)
		c.JSON(http.StatusOK, gin.H{
			"msg": msg,
		})
	})

	// 查（列表）
	router.GET("/person", func(c *gin.Context) {
		var person []Person
		err := db.Table("person").Order("id desc").Find(&person).Error

		if err != nil {
			c.JSON(http.StatusCreated, gin.H{
				"person": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"person": person,
		})
	})

	// 查（单条）
	router.GET("/person/:id", func(c *gin.Context) {
		id := c.Param("id")

		var person Person
		err := db.Table("person").Where("id = ?", id).Find(&person).Error

		// 数据不存在
		// sql: no rows in result set
		if err != nil {
			c.JSON(http.StatusCreated, gin.H{
				"person": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"person": person,
		})
	})

	// 改
	router.PUT("/person/:id", func(c *gin.Context) {
		cid := c.Param("id")
		id, err := strconv.Atoi(cid)
		person := Person{Id: id}
		err = c.Bind(&person)

		if err != nil {
			panic(err)
		}

		stmt, err := db.DB().Prepare("UPDATE person SET first_name=?, last_name=? WHERE id=?")

		err = stmt.Close()
		if err != nil {
			panic(err)
		}

		rs, err := stmt.Exec(person.FirstName, person.LastName, person.Id)
		if err != nil {
			panic(err)
		}

		ra, err := rs.RowsAffected()
		if err != nil {
			panic(err)
		}

		msg := fmt.Sprintf("Update person %d successful %d", person.Id, ra)
		c.JSON(http.StatusOK, gin.H{
			"msg": msg,
		})
	})

	// 删
	router.DELETE("/person/:id", func(c *gin.Context) {
		cid := c.Param("id")

		id, err := strconv.Atoi(cid)
		if err != nil {
			panic(err)
		}

		rs, err := db.DB().Exec("DELETE FROM person WHERE id=?", id)
		if err != nil {
			panic(err)
		}

		ra, err := rs.RowsAffected()
		if err != nil {
			panic(err)
		}

		msg := fmt.Sprintf("Delete person %d successful %d", id, ra)
		c.JSON(http.StatusOK, gin.H{
			"msg": msg,
		})
	})

	port := "8888"

	err = router.Run(":" + port)
	if err != nil {
		panic(err)
	}

	// http.Server 内置的 Shutdown 方法来实现优雅的关闭服务
	go _shutdown(port, router)
}

func _shutdown(port string, router *gin.Engine) {
	server := http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		panic("Server Shutdown: " + err.Error())
	}

	var c *gin.Context

	c.JSON(http.StatusOK, gin.H{
		"msg": "Server exiting",
	})
}
