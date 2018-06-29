package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

func init() {
	var err error
	db, err = sqlx.Connect("mysql", "root:3411@tcp(127.0.0.1:3306)/pshdata")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	r := gin.Default()
	r.GET("/ping/", PingMe)
	r.GET("/states/", GetStates)

	r.Run()
	//fmt.Printf("hello, world\n")
}

//PingMe check db state
func PingMe(c *gin.Context) {
	var res ValidateResult
	res.State = -1000
	if err := db.Ping(); err != nil {
		res.ErrCode = -1
		res.Message = err.Error()
	} else {
		res.Message = "Ping OK"
	}
	c.JSON(200, res)

}

//GetStates read ClientState's from db
func GetStates(c *gin.Context) {
	if list, err := loadStates(); err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, list)
	}
}
