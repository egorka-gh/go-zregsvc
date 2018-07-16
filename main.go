package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/dchest/captcha"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

var db *sqlx.DB

func init() {
	var err error
	viper.SetDefault("Port", "8080")
	viper.SetDefault("ConnectionString", "root:3411@tcp(127.0.0.1:3306)/pshdata")
	//TODO remove in prod
	viper.SetDefault("cors", "http://localhost:4200")
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	err = viper.ReadInConfig()
	if err != nil {
		log.Print(err)
		log.Print("Using defaults")
	}

	db, err = sqlx.Connect("mysql", viper.GetString("ConnectionString"))
	if err != nil {
		log.Fatal(err)
	}

	//gin.SetMode(gin.ReleaseMode)
}

func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = strings.Split(viper.GetString("cors"), ",")
	//config.AllowOrigins = []string{"http://localhost:4200"}
	//config.AllowAllOrigins = true

	r.Use(cors.New(config))
	r.GET("/ping/", PingMe)
	r.GET("/states/", GetStates)
	r.GET("/captcha/:image", GetCaptcha)
	r.POST("/validate/", CardValidate)

	r.Run(":" + viper.GetString("Port"))
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
		res.Captcha = captcha.New()
		res.CaptchaState = 1
	}
	//log.Println("PingMe")
	c.JSON(200, res)
}

//GetCaptcha check db state
func GetCaptcha(c *gin.Context) {
	//gin.WrapH(captcha.Server(200, 80))
	captcha.Server(210, 70).ServeHTTP(c.Writer, c.Request)
}

//GetStates read ClientState's from db
func GetStates(c *gin.Context) {
	if list, err := loadStates(); err != nil {
		c.AbortWithStatus(404)
		log.Print(err)
	} else {
		c.JSON(200, list)
	}
}

//CardValidate check client
func CardValidate(c *gin.Context) {
	var res ValidateResult

	if err := c.ShouldBindJSON(&res); err == nil {
		c.JSON(http.StatusOK, validateCard(res.Card))
	} else {
		res.ErrCode = -1
		res.Message = err.Error()
		c.JSON(http.StatusOK, res)
	}
}
