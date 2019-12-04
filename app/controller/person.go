package controller

import (
	"fmt"
	"gin/internal/databases/mysql"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Person struct {
	Id        int    `gorm:"column:id; PRIMARY_KEY ;AUTO_INCREMENT" json:"id" form:"id"`
	FirstName string `gorm:"column:first_name" json:"first_name" form:"first_name"`
	LastName  string `gorm:"column:last_name" json:"last_name" form:"last_name"`
}

// 增加
func AddPerson(c *gin.Context) {
	firstName := c.Request.FormValue("first_name")
	lastName := c.Request.FormValue("last_name")

	db := mysql.Client()

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
}

// 查（列表）
func ListPerson(c *gin.Context) {
	var person []Person

	db := mysql.Client()

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
}

// 查（单条）
func FindPerson(c *gin.Context) {
	id := c.Param("id")

	db := mysql.Client()

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
}

// 改
func EditPerson(c *gin.Context) {
	cid := c.Param("id")
	id, err := strconv.Atoi(cid)
	person := Person{Id: id}
	err = c.Bind(&person)

	if err != nil {
		panic(err)
	}

	db := mysql.Client()

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
}

// 删
func DelPerson(c *gin.Context) {
	cid := c.Param("id")

	id, err := strconv.Atoi(cid)
	if err != nil {
		panic(err)
	}

	db := mysql.Client()

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
}

