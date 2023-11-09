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

	provider.RedirectUri = domain + path.Callback

	return NewNinshoGin[T](r, provider, idp, domain, &path)
}

func (_ninsho *NinshoGin[T]) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		user, err := LoadUser[ninsho.LINE_USER](c)

		if user == nil || err != nil {
			c.Redirect(http.StatusTemporaryRedirect, _ninsho.Path.Unauthorized)
			c.AbortWithStatus(401)
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

		err = SaveUser(jwt, c)

		if err != nil {
			panic(err)
		}

		c.Redirect(http.StatusTemporaryRedirect, ninsho.Path.AfterAuth)
	})
}
