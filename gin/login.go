package gin_line_login

import (
	"net/http"

	core "line_login_core"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const userId = "userId"

type LineLogin struct {
	UnauthorizedPath string
	CallbackPath     string
	AfterAuthPath    string
	Domain           string
	LineLoginSession *core.Session
}

func NewLineLogin(r *gin.RouterGroup, provider *core.Provider, domain, unauthorized, callback, afterAuth string) (*LineLogin, error) {
	session, err := core.NewSession(provider)
	if err != nil {
		return nil, err
	}

	lineLogin := LineLogin{
		UnauthorizedPath: unauthorized,
		CallbackPath:     callback,
		AfterAuthPath:    afterAuth,
		Domain:           domain,
	}

	lineLogin.LineLoginSession = &session

	lineLogin.Callback(r)

	return &lineLogin, nil
}

func NewLineLoginWithEnvironment(r *gin.RouterGroup, unauthorized, callback, afterAuth string) (*LineLogin, error) {
	var provider core.Provider

	domain := os.Getenv("LINE_LOGIN_BASE")

	provider.ClientID = os.Getenv("LINE_LOGIN_CLIENT_ID")
	provider.ClientSecret = os.Getenv("LINE_LOGIN_CLIENT_SECRET")
	provider.RedirectUri = domain + callback

	return NewLineLogin(r, &provider, domain, unauthorized, callback, afterAuth)
}

func DefaultLineLogin(r *gin.RouterGroup) (*LineLogin, error) {
	return NewLineLoginWithEnvironment(r, "/unauthorized", "/callback", "/")
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

func (lineLogin *LineLogin) Callback(r *gin.RouterGroup) {
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
