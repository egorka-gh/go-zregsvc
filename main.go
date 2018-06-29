package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

//ClientState clent states & frontend messages
type ClientState struct {
	ID         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	WebComment string `json:"сomment" db:"web_comment"`
}

//ValidateResult response 4 frontend
type ValidateResult struct {
	Program int    `json:"program" db:"program"`
	Card    string `json:"card" db:"card"`
	State   int    `json:"state" db:"state"`
	ErrCode int    `json:"err" db:"err"`
	Message string `json:"message" db:"message"`
}

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
	if err := db.Ping(); err != nil {
		c.JSON(200, gin.H{
			"message": err.Error,
		})
	} else {
		c.JSON(200, gin.H{
			"message": "Ping OK",
		})
	}
}

//GetStates read ClientState's from db
func GetStates(c *gin.Context) {
	var list []ClientState
	var ssql = "SELECT cs.id, ifnull(cs.name,'') name, ifnull(csm.web_comment,'') web_comment FROM client_state cs LEFT OUTER JOIN client_state_msg csm ON cs.id = csm.id WHERE cs.id!=0 ORDER BY cs.id"
	if err := db.Select(&list, ssql); err != nil {
		c.AbortWithStatus(404)
		fmt.Println(err)
	} else {
		c.JSON(200, list)
	}

}
