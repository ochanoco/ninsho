package gin_line_login

import (
	"net/http"

	"line_login_core"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const userId = "userId"

type LineLogin struct {
	UnauthorizedPath string
	CallbackPath     string
	AfterAuthPath    string
	LineLoginSession *line_login_core.Session
}

func NewLineLogin(r *gin.Engine, unauthorized, callback, afterAuth string) (*LineLogin, error) {
	var provider line_login_core.Provider

	provider.ClientID = os.Getenv("CLIENT_ID")
	provider.ClientSecret = os.Getenv("CLIENT_SECRET")
	provider.RedirectUri = os.Getenv("REDIRECT_URI") + callback

	session, err := line_login_core.NewSession(&provider)
	if err != nil {
		return nil, err
	}

	lineLogin := LineLogin{
		UnauthorizedPath: unauthorized,
		CallbackPath:     callback,
		AfterAuthPath:    afterAuth,
	}

	lineLogin.LineLoginSession = &session

	lineLogin.Callback(r)

	return &lineLogin, nil
}

func DefaultLineLogin(r *gin.Engine) (*LineLogin, error) {
	return NewLineLogin(r, "/unauthorized", "/callback", "/")
}

func (lineLogin *LineLogin) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get(userId) == nil {
			c.Redirect(http.StatusTemporaryRedirect, lineLogin.UnauthorizedPath)
			c.Abort()
		}

		c.Next()
	}
}

func (lineLogin *LineLogin) Login(c *gin.Context) {
	url := lineLogin.LineLoginSession.AuthURL()

	c.Redirect(http.StatusTemporaryRedirect, url)
	c.Abort()
}

func (lineLogin *LineLogin) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(userId)
	session.Save()
}

func (lineLogin *LineLogin) Callback(r *gin.Engine) {
	r.GET(lineLogin.CallbackPath, func(c *gin.Context) {
		code := c.Query("code")

		jwt, err := lineLogin.LineLoginSession.GetUser(code)
		if err != nil {
			panic(err)
		}

		session := sessions.Default(c)

		session.Set(userId, jwt.Sub)
		session.Save()

		c.Redirect(http.StatusTemporaryRedirect, lineLogin.AfterAuthPath)
	})
}
