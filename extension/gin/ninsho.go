package gin_ninsho

import (
	"net/http"
	"os"

	"ninsho"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const userId = "userId"

type NinshoGin struct {
	UnauthorizedPath string
	CallbackPath     string
	AfterAuthPath    string
	Domain           string
	Ninsho           *ninsho.Ninsho
}

func NewNinshoGin(r *gin.RouterGroup, provider *ninsho.Provider, domain, unauthorized, callback, afterAuth string) (*NinshoGin, error) {
	session, err := ninsho.NewNinsho(provider)
	if err != nil {
		return nil, err
	}

	ninshoGin := NinshoGin{
		UnauthorizedPath: unauthorized,
		CallbackPath:     callback,
		AfterAuthPath:    afterAuth,
		Domain:           domain,
	}

	ninshoGin.Ninsho = &session

	ninshoGin.Callback(r)

	return &ninshoGin, nil
}

func NewNinshoGinFromEnv(r *gin.RouterGroup, unauthorized, callback, afterAuth string) (*NinshoGin, error) {
	var provider ninsho.Provider

	domain := os.Getenv("LINE_LOGIN_BASE")
	provider.ClientID = os.Getenv("LINE_LOGIN_CLIENT_ID")
	provider.ClientSecret = os.Getenv("LINE_LOGIN_CLIENT_SECRET")
	provider.RedirectUri = domain + callback

	return NewNinshoGin(r, &provider, domain, unauthorized, callback, afterAuth)
}

func DefaultNinshoGin(r *gin.RouterGroup) (*NinshoGin, error) {
	return NewNinshoGinFromEnv(r, "/unauthorized", "/callback", "/")
}

func (ninsho *NinshoGin) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get(userId) == nil {
			c.Redirect(http.StatusTemporaryRedirect, ninsho.UnauthorizedPath)
			c.Abort()
		}

		c.Next()
	}
}

func (ninsho *NinshoGin) Login(c *gin.Context) {
	url := ninsho.Ninsho.AuthURL()

	c.Redirect(http.StatusTemporaryRedirect, url)
	c.Abort()
}

func (ninsho *NinshoGin) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete(userId)
	session.Save()
}

func (ninsho *NinshoGin) Callback(r *gin.RouterGroup) {
	r.GET(ninsho.CallbackPath, func(c *gin.Context) {
		code := c.Query("code")

		jwt, err := ninsho.Ninsho.GetUser(code)
		if err != nil {
			panic(err)
		}

		session := sessions.Default(c)

		session.Set(userId, jwt.Sub)
		session.Save()

		c.Redirect(http.StatusTemporaryRedirect, ninsho.AfterAuthPath)
	})
}
