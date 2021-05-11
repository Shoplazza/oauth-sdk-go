package oauth2

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
)

var DefaultRequestPath = "/auth/shoplazza"
var DefaultCallbackPath = "/auth/shoplazza/callback"

type GinMiddleware struct {
	OauthConfig  *Config
	RequestPath  string
	CallbackPath string

	// ProviderIgnoreState bool
	// 自定义 request 处理函数
	RequestFunc func(c *gin.Context)
	// 自定义 callback 处理函数
	CallbackFunc func(c *gin.Context)
}

func NewGinMiddleware(oauthConfig *Config, options ...OauthGinOptionFunc) *GinMiddleware {
	gm := &GinMiddleware{
		RequestPath:  DefaultRequestPath,
		CallbackPath: DefaultCallbackPath,
		OauthConfig:  oauthConfig,
	}

	for _, option := range options {
		option(gm)
	}
	return gm
}

func (gm *GinMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		switch c.Request.URL.Path {
		case gm.RequestPath:
			if gm.RequestFunc != nil {
				gm.RequestFunc(c)
			} else {
				gm.ginOauthRequest(c)
			}
		case gm.CallbackPath:
			if gm.CallbackFunc != nil {
				gm.CallbackFunc(c)
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
	if params == nil || !gm.OauthConfig.ValidShop(params.Get("shop")) {
		c.String(400, "OAuth endpoint is not a myshoplazza site.")
		return
	}

	var opts []AuthCodeOption
	// TODO state
	//if !gm.ProviderIgnoreState {
	//	state := GetRandomString(48)
	//	opts = append(opts, SetAuthURLParam("state", state))
	//	session := sessions.Default(c)
	//	session.Set("auth2.state", state)
	//}
	c.Redirect(302, gm.OauthConfig.AuthCodeURL(params.Get("shop"), opts...))
}

func (gm *GinMiddleware) ginOauthCallback(c *gin.Context) {
	// TODO valid state => message : CSRF detected

	params := gm.getParams(c)
	if params == nil {
		c.String(400, "Invalid callback")
		c.Abort()
		return
	}

	shop := params.Get("shop")
	if !gm.OauthConfig.ValidShop(shop) {
		c.String(400, "OAuth endpoint is not a myshoplazza site.")
		c.Abort()
		return
	}

	if !gm.signatureValid(params) {
		c.String(400, "Signature does not match, it may have been tampered with.")
		c.Abort()
		return
	}

	token, err := gm.OauthConfig.Exchange(context.Background(), shop, params.Get("code"))
	if err != nil {
		c.String(400, err.Error())
		c.Abort()
		return
	}
	c.Set("oauth2.token", token)
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

func (gm *GinMiddleware) signatureValid(params url.Values) bool {
	v := params.Get("hmac")
	params.Del("hmac")

	hm := hmac.New(sha256.New, []byte(gm.OauthConfig.ClientSecret))
	hm.Write([]byte(params.Encode()))
	signature := hex.EncodeToString(hm.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(v))
}

type OauthGinOptionFunc func(*GinMiddleware)

func SetRequestPath(path string) OauthGinOptionFunc {
	return func(gm *GinMiddleware) {
		gm.RequestPath = path
	}
}

func SetCallbackPath(path string) OauthGinOptionFunc {
	return func(gm *GinMiddleware) {
		gm.CallbackPath = path
	}
}

func GetRandomString(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}
