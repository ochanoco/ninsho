package main

import (
	"gin_ninsho"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	ninsho, err := gin_ninsho.DefaultNinshoGin(&r.RouterGroup)
	if err != nil {
		panic(err)
	}

	r.GET("/", ninsho.AuthMiddleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "loggined!"})
	})

	r.GET("/login", func(c *gin.Context) {
		ninsho.Login(c)
	})

	r.GET("/logout", func(c *gin.Context) {
		ninsho.Logout(c)
		c.JSON(200, gin.H{"message": "logout"})
	})

	r.GET("/unauthorized", func(c *gin.Context) {
		c.JSON(401, gin.H{"message": "unauthorized"})
	})
	r.Run()
}
