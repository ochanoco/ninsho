package gin_ninsho

import (
	"encoding/json"
	"net/http"
	"os"

	"ninsho"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type NinshoGin[T any] struct {
	UnauthorizedPath string
	CallbackPath     string
	AfterAuthPath    string
	Domain           string
	Ninsho           *ninsho.Ninsho[T]
}

func NewNinshoGin[T any](r *gin.RouterGroup, provider *ninsho.Provider, idp *ninsho.IdP[T], domain, unauthorized, callback, afterAuth string) (*NinshoGin[T], error) {
	session, err := ninsho.NewNinsho[T](provider, idp)
	if err != nil {
		return nil, err
	}

	ninshoGin := NinshoGin[T]{
		UnauthorizedPath: unauthorized,
		CallbackPath:     callback,
		AfterAuthPath:    afterAuth,
		Domain:           domain,
	}

	ninshoGin.Ninsho = &session

	ninshoGin.Callback(r)

	return &ninshoGin, nil
}

func NewNinshoGinFromEnv[T any](r *gin.RouterGroup, idp *ninsho.IdP[T], unauthorized, callback, afterAuth string) (*NinshoGin[T], error) {
	var provider ninsho.Provider

	domain := os.Getenv("NINSHO_BASE")
	provider.ClientID = os.Getenv("NINSHO_CLIENT_ID")
	provider.ClientSecret = os.Getenv("NINSHO_CLIENT_SECRET")
	provider.RedirectUri = domain + callback

	return NewNinshoGin[T](r, &provider, idp, domain, unauthorized, callback, afterAuth)
}

func DefaultNinshoGin[T any](r *gin.RouterGroup, idp *ninsho.IdP[T]) (*NinshoGin[T], error) {
	return NewNinshoGinFromEnv[T](r, idp, "/unauthorized", "/callback", "/")
}

func (ninsho *NinshoGin[T]) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get("user") == nil {
			c.Redirect(http.StatusTemporaryRedirect, ninsho.UnauthorizedPath)
			c.Abort()
		}

		c.Next()
	}
}

func (ninsho *NinshoGin[T]) Login(c *gin.Context) {
	url := ninsho.Ninsho.GetAuthURL()

	c.Redirect(http.StatusTemporaryRedirect, url)
	c.Abort()
}

func (ninsho *NinshoGin[T]) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user")
	session.Save()
}

func (ninsho *NinshoGin[T]) Callback(r *gin.RouterGroup) {
	r.GET(ninsho.CallbackPath, func(c *gin.Context) {
		code := c.Query("code")

		jwt, err := ninsho.Ninsho.Auth(code)
		if err != nil {
			panic(err)
		}

		user, err := json.Marshal(jwt)
		if err != nil {
			panic(err)
		}

		session := sessions.Default(c)

		session.Set("user", user)
		session.Save()

		c.Redirect(http.StatusTemporaryRedirect, ninsho.AfterAuthPath)
	})
}