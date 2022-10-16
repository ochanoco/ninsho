package gin_line_login

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"line_login_core"
	"os"
)

const ID_TOKEN = "id_token"

type LineLogin struct {
	UnauthorizedPath string
	CallbackPath     string
	LineLoginSession *line_login_core.Session
}

func NewLineLogin(r *gin.Engine, unauthorized, callback string) (*LineLogin, error) {
	var provider line_login_core.Provider

	provider.ClientID = os.Getenv("CLIENT_ID")
	provider.ClientSecret = os.Getenv("TOKEN")
	provider.RedirectURL = "http://127.0.0.1:8080/callback"

	session, err := line_login_core.NewSession(&provider)
	if err != nil {
		return nil, err
	}

	lineLogin := LineLogin{
		UnauthorizedPath: unauthorized,
		CallbackPath:     callback,
	}

	lineLogin.LineLoginSession = &session

	lineLogin.Callback(r)

	return &lineLogin, nil
}

func DefaultLineLogin(r *gin.Engine) (*LineLogin, error) {
	return NewLineLogin(r, "/unauthorized", "/callback")
}

func (lineLogin *LineLogin) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get(ID_TOKEN) == nil {
			c.Redirect(http.StatusFound, lineLogin.UnauthorizedPath)
			c.Abort()
		}

		c.Next()
	}
}

func (lineLogin *LineLogin) Login(c *gin.Context) {
	url := lineLogin.LineLoginSession.AuthURL()

	c.Redirect(http.StatusMovedPermanently, url)
	c.Abort()
}

func (lineLogin *LineLogin) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(ID_TOKEN)
	session.Save()
}

func (lineLogin *LineLogin) Callback(r *gin.Engine) {
	r.GET(lineLogin.CallbackPath, func(c *gin.Context) {
		code := c.Query("code")

		user, err := lineLogin.LineLoginSession.GetUser(code)
		if err != nil {
			panic(err)
		}

		c.JSON(200, gin.H{ID_TOKEN: user.IdToken})
	})
}
