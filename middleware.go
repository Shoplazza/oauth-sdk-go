package common_oauth2

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

var DefaultRequestPath = "/auth/shoplazza"
var DefaultCallbackPath = "/auth/shoplazza/callback"
var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(32), securecookie.GenerateRandomKey(32))

type GinMiddleware struct {
	oauthConfig  *Config
	requestPath  string
	callbackPath string
	ignoreState  bool

	// 自定义 request 处理函数
	requestFunc func(c *gin.Context)
	// 自定义 callback 处理函数
	callbackFunc func(c *gin.Context)
	// access-token 处理函数
	accessTokenHandlerFunc func(c *gin.Context, shop string, token *Token)
}

func NewGinMiddleware(oauthConfig *Config) *GinMiddleware {
	gm := &GinMiddleware{
		requestPath:  DefaultRequestPath,
		callbackPath: DefaultCallbackPath,
		oauthConfig:  oauthConfig,
	}
	return gm
}

func (gm *GinMiddleware) IgnoreState(ignoreState bool) {
	gm.ignoreState = ignoreState
}

func (gm *GinMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.URL.Path {
		case gm.requestPath:
			if gm.requestFunc != nil {
				gm.requestFunc(c)
			} else {
				gm.ginOauthRequest(c)
			}
		case gm.callbackPath:
			if gm.callbackFunc != nil {
				gm.callbackFunc(c)
			} else {
				gm.ginOauthCallback(c)
			}
		default:
			c.Next()
		}
	}
}

func (gm *GinMiddleware) ginOauthRequest(c *gin.Context) {
	params := gm.getParams(c)
	if params == nil || !gm.oauthConfig.ValidShop(params.Get("shop")) {
		c.String(400, "OAuth endpoint is not a myshoplazza site.")
		return
	}

	var opts []AuthCodeOption

	if gm.ignoreState == false {
		state := GetRandomString(48)
		opts = append(opts, SetAuthURLParam("state", state))

		session, _ := store.Get(c.Request, "state-session")
		session.Values["state"] = state
		session.Save(c.Request, c.Writer)
	}

	c.Redirect(302, gm.oauthConfig.AuthCodeURL(params.Get("shop"), opts...))
}

func (gm *GinMiddleware) ginOauthCallback(c *gin.Context) {
	params := gm.getParams(c)
	if params == nil {
		c.String(400, "Invalid callback")
		c.Abort()
		return
	}

	shop := params.Get("shop")
	if !gm.oauthConfig.ValidShop(shop) {
		c.String(400, "OAuth endpoint is not a myshoplazza site.")
		c.Abort()
		return
	}

	if gm.ignoreState == false {
		stateFromParam := params.Get("state")
		session, _ := store.Get(c.Request, "state-session")
		stateFromSession := session.Values["state"]
		if stateFromSession == nil {
			c.String(400, "State does not exist in the session.")
			c.Abort()
			return
		}
		if stateFromParam != stateFromSession.(string) {
			c.String(400, "State does not match.")
			c.Abort()
			return
		}
	}

	if !gm.oauthConfig.SignatureValid(params) {
		c.String(400, "Signature does not match, it may have been tampered with.")
		c.Abort()
		return
	}

	token, err := gm.oauthConfig.Exchange(context.Background(), shop, params.Get("code"))
	if err != nil {
		c.String(400, err.Error())
		c.Abort()
		return
	}
	c.Set("oauth2.token", token)

	if gm.accessTokenHandlerFunc == nil {
		c.Next()
	} else {
		gm.accessTokenHandlerFunc(c, shop, token)
	}
}

func (gm *GinMiddleware) getParams(c *gin.Context) url.Values {
	query := strings.Split(c.Request.RequestURI, "?")
	if len(query) != 2 {
		return nil
	}

	params, err := url.ParseQuery(query[1])
	if err != nil {
		return nil
	}
	return params
}

func (gm *GinMiddleware) SetRequestPath(path string) {
	gm.requestPath = path
}

func (gm *GinMiddleware) SetCallbackPath(path string) {
	gm.callbackPath = path
}

func (gm *GinMiddleware) SetRequestFunc(fn func(c *gin.Context)) {
	gm.requestFunc = fn
}

func (gm *GinMiddleware) SetCallbackFunc(fn func(c *gin.Context)) {
	gm.callbackFunc = fn
}

func (gm *GinMiddleware) SetAccessTokenHandlerFunc(fn func(c *gin.Context, shop string, token *Token)) {
	gm.accessTokenHandlerFunc = fn
}

func GetRandomString(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}
