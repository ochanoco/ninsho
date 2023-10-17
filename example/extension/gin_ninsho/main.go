package main

import (
	"os"

	"gin_ninsho"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/ochanoco/ninsho"
)

var DOMAIN = os.Getenv("NINSHO_BASE")
var CLIENT_ID = os.Getenv("NINSHO_CLIENT_ID")
var CLIENT_SECRET = os.Getenv("NINSHO_CLIENT_SECRET")

func main() {
	r := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	provider := ninsho.Provider{
		ClientID:     CLIENT_ID,
		ClientSecret: CLIENT_SECRET,
	}

	_ninsho, err := gin_ninsho.DefaultNinshoGin(&r.RouterGroup, provider, ninsho.LINE_LOGIN, DOMAIN, CLIENT_ID, CLIENT_SECRET)
	if err != nil {
		panic(err)
	}

	r.GET("/", _ninsho.AuthMiddleware(), func(c *gin.Context) {
		jwt, err := gin_ninsho.GetUser[ninsho.LINE_USER](c)
		if err != nil {
			panic(err)
		}

		c.JSON(200, gin.H{"message": "loggined!", "user": jwt.Sub})
	})

	r.GET("/login", func(c *gin.Context) {
		_ninsho.Login(c)
	})

	r.GET("/logout", func(c *gin.Context) {
		_ninsho.Logout(c)
		c.JSON(200, gin.H{"message": "logout"})
	})

	r.GET("/unauthorized", func(c *gin.Context) {
		c.JSON(401, gin.H{"message": "unauthorized"})
	})
	r.Run()
}
