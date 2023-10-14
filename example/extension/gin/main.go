package main

import (
	"ninsho"
	gin_ninsho "ninsho/extension/gin"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	_ninsho, err := gin_ninsho.DefaultNinshoGin(&r.RouterGroup, &ninsho.LINE_LOGIN)
	if err != nil {
		panic(err)
	}

	r.GET("/", _ninsho.AuthMiddleware(), func(c *gin.Context) {
		jwt, err := gin_ninsho.GetUser[ninsho.LINE_JWT](c)
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
