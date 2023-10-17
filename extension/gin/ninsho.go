package gin

import (
	"encoding/json"
	"net/http"

	"github.com/ochanoco/ninsho"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type NinshoGinPath struct {
	Unauthorized string
	Callback     string
	AfterAuth    string
}

type NinshoGin[T any] struct {
	Path   *NinshoGinPath
	Domain string
	Ninsho *ninsho.Ninsho[T]
}

func NewNinshoGin[T any](r *gin.RouterGroup, provider *ninsho.Provider, idp *ninsho.IdP[T], domain string, path *NinshoGinPath) (*NinshoGin[T], error) {
	n, err := ninsho.NewNinsho[T](provider, idp)
	if err != nil {
		return nil, err
	}

	ninshoGin := NinshoGin[T]{
		Domain: domain,
		Path:   path,
		Ninsho: &n,
	}

	ninshoGin.Callback(r)

	return &ninshoGin, nil
}

func DefaultNinshoGin[T any](r *gin.RouterGroup, provider *ninsho.Provider, idp *ninsho.IdP[T], domain string) (*NinshoGin[T], error) {
	path := NinshoGinPath{
		Unauthorized: "/unauthorized",
		Callback:     "/callback",
		AfterAuth:    "/",
	}
	return NewNinshoGin[T](r, provider, idp, domain, &path)
}

func (ninsho *NinshoGin[T]) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		if session.Get("user") == nil {
			c.Redirect(http.StatusTemporaryRedirect, ninsho.Path.Unauthorized)
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
	r.GET(ninsho.Path.Callback, func(c *gin.Context) {
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

		c.Redirect(http.StatusTemporaryRedirect, ninsho.Path.AfterAuth)
	})
}
