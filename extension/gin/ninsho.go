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
	BaseUrl  string
	Provider *ninsho.Provider
	IdP      *ninsho.IdP[T]
}

func NewNinshoGin[T any](r *gin.RouterGroup, provider *ninsho.Provider, idp *ninsho.IdP[T], baseUrl string, path *NinshoGinPath) (*NinshoGin[T], error) {
	ninshoGin := NinshoGin[T]{
		BaseUrl:  baseUrl,
		Path:     path,
		Provider: provider,
		IdP:      idp,
	}

	ninshoGin.Callback(r)

	return &ninshoGin, nil
}

func DefaultNinshoGin[T any](r *gin.RouterGroup, provider *ninsho.Provider, idp *ninsho.IdP[T], baseUri string) (*NinshoGin[T], error) {
	path := NinshoGinPath{
		Unauthorized: "/unauthorized",
		Callback:     "/callback",
		AfterAuth:    "/",
	}

	provider.RedirectUri = baseUri + path.Callback

	return NewNinshoGin(r, provider, idp, baseUri, &path)
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
	n, err := ninsho.InitNinsho(_ninsho.Provider, _ninsho.IdP)
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

		n := ninsho.NewNinsho(_ninsho.Provider, _ninsho.IdP, stateSession)

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
