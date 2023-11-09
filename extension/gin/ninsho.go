package gin

import (
	"net/http"

	"github.com/ochanoco/ninsho"

	// "ninsho"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type NinshoGinPath struct {
	Unauthorized string
	Callback     string
	AfterAuth    string
}

type NinshoGin[T any] struct {
	Path     *NinshoGinPath
	Domain   string
	Provider *ninsho.Provider
	IdP      *ninsho.IdP[T]
}

func NewNinshoGin[T any](r *gin.RouterGroup, provider *ninsho.Provider, idp *ninsho.IdP[T], domain string, path *NinshoGinPath) (*NinshoGin[T], error) {
	ninshoGin := NinshoGin[T]{
		Domain:   domain,
		Path:     path,
		Provider: provider,
		IdP:      idp,
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

	provider.RedirectUri = domain + path.Callback

	return NewNinshoGin[T](r, provider, idp, domain, &path)
}

func (_ninsho *NinshoGin[T]) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := LoadUser[ninsho.LINE_USER](c)

		if user == nil || err != nil {
			c.Redirect(http.StatusTemporaryRedirect, _ninsho.Path.Unauthorized)
			c.AbortWithStatus(401)
			return
		}

		c.Next()
	}
}

func (_ninsho *NinshoGin[T]) Login(c *gin.Context) {
	n, err := ninsho.InitNinsho[T](_ninsho.Provider, _ninsho.IdP)
	if err != nil {
		panic(err)
	}

	url, _, err := n.MakeAuthURL()
	if err != nil {
		panic(err)
	}

	SaveLoginingSession(*n.Session, c)

	c.Redirect(http.StatusTemporaryRedirect, url)
	c.Abort()
}

func (ninsho *NinshoGin[T]) Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("user")
	session.Save()
}

func (_ninsho *NinshoGin[T]) Callback(r *gin.RouterGroup) {
	r.GET(_ninsho.Path.Callback, func(c *gin.Context) {
		code := c.Query("code")
		state := c.Query("state")

		stateSession, err := LoadLoginingSession(c)
		if err != nil {
			panic(err)
		}

		n := ninsho.NewNinsho[T](_ninsho.Provider, _ninsho.IdP, stateSession)

		jwt, err := n.Auth(code, state)
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, _ninsho.Path.Unauthorized)
			c.AbortWithStatus(401)
			return
		}

		err = SaveUser(jwt, c)

		if err != nil {
			panic(err)
		}

		c.Redirect(http.StatusTemporaryRedirect, _ninsho.Path.AfterAuth)
	})
}
