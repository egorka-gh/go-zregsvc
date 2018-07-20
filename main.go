package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dchest/captcha"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

var db *sqlx.DB
var captchaStore TagedStore

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

	captchaStore = NewTagedStore(100, 60*time.Minute)
	captcha.SetCustomStore(captchaStore)

	//TODO add in prod
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
	r.POST("/register/", CardRegister)

	r.Run(":" + viper.GetString("Port"))
	//fmt.Printf("hello, world\n")
}

//PingMe check db state
func PingMe(c *gin.Context) {
	var res ValidateResult
	res.State = 0
	if err := db.Ping(); err != nil {
		res.ErrCode = errDATABASE
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

func validateCaptcha(res *ValidateResult) bool {
	//check captcha
	if len(res.Captcha) == 0 {
		//err captha id empty
		res.ErrCode = errWrongCAPTCHA
		res.CaptchaState = -10
		return false
	}
	card, ok := captchaStore.GetTag(res.Captcha)
	if !ok {
		//err captha id invalid
		res.ErrCode = errWrongCAPTCHA
		res.CaptchaState = -10
		return false
	}
	//chek if captha not taged vs card or card is same
	if len(card) != 0 && card != res.Card {
		//err attempt to use solved capthca vs onother card
		res.ErrCode = errWrongCAPTCHA
		res.CaptchaState = -100
		return false
	}

	if !captcha.VerifyString(res.Captcha, res.CaptchaSolution) {
		//err wrong  captha solution
		res.ErrCode = errWrongCAPTCHA
		res.CaptchaState = -1
		return false
	}
	res.ErrCode = 0
	res.CaptchaState = 100
	return true
}

//CardValidate check client
func CardValidate(c *gin.Context) {
	var res ValidateResult

	if err := c.ShouldBindJSON(&res); err == nil {
		//check captcha
		if !validateCaptcha(&res) {
			c.JSON(http.StatusOK, res)
			return
		}

		//check card in db
		res = validateCard(res)
		if res.ErrCode == 0 {
			//lock capthca vc card
			if !captchaStore.SetTag(res.Captcha, res.Card) {
				// some thread lock capthca ???
				res.ErrCode = errWrongCAPTCHA
				res.CaptchaState = -100
				res.Program = 0
				res.State = 0
			}
		}
		c.JSON(http.StatusOK, res)
	} else {
		res.ErrCode = errDATABASE
		res.Message = err.Error()
		c.JSON(http.StatusOK, res)
	}
}

//CardRegister register client
func CardRegister(c *gin.Context) {
	var dto RegisterDTO
	var res ValidateResult

	if err := c.ShouldBindJSON(&dto); err == nil {
		res = dto.Result
		//check captcha
		if !validateCaptcha(&res) {
			c.JSON(http.StatusOK, res)
			return
		}

		captchaStore.Del(res.Captcha)

		if res.Program == 0 || res.Card != dto.Client.Card {
			res.ErrCode = errDATABASE
			res.Message = "Application flow error"
			c.JSON(http.StatusOK, res)
		}
		//set client data

		c.JSON(http.StatusOK, registerCard(&dto))
	} else {
		res.ErrCode = errDATABASE
		res.Message = err.Error()
		c.JSON(http.StatusOK, res)
	}
}
