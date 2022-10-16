package gin_line_login

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LineLogin struct {
	UnauthorizedPath string
}

func NewLineLogin(unauthorized string) *LineLogin {
	return &LineLogin{
		unauthorized,
	}
}

func DefaultLineLogin() *LineLogin {
	return NewLineLogin("/unauthorized")
}

func (lineLogin *LineLogin) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get("id_token") == nil {
			c.Redirect(http.StatusFound, lineLogin.UnauthorizedPath)
			c.Abort()
		}

		c.Next()
	}
}

func (lineLogin *LineLogin) Login(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("id_token", "123")
	session.Save()
}

func (lineLogin *LineLogin) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("id_token")
	session.Save()
}
