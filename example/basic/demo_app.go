package main

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	co "github.com/Shoplazza/oauth-sdk-go"
	"github.com/Shoplazza/oauth-sdk-go/shoplazza"
)

type Token struct {
	Id          int       `json:"id"`
	AccessToken string    `json:"access_token"`
	Shop        string    `json:"shop"`
	StoreId     string    `json:"store_id"`
	ExpiredAt   time.Time `json:"expired_at"`
}

func main() {
	oauth := &co.Config{
		ClientID:     "WVpHNkENUL9CBbDjGO_po9tfG02XW2Z-X54M4LObfDs", // your app client id
		ClientSecret: "DdJNhsopKxAWDHjqjI1rpQZW17Fp6GXrHhC0IgwdXag", // your app client secret
		Endpoint:     shoplazza.Endpoint,
		RedirectURI:  "https://6f1e-123-58-221-57.ngrok.io/oauth_sdk/redirect_uri", // your app redirect uri, you need to replace the 6f1e-123-58-221-57.ngrok.io with your service domain
		Scopes:       []string{"read_shop", "write_shop"},                          // []string{"scope1", "scope2"}
	}

	r := gin.New()

	// app uri endpoint
	// if your redirect uri is 'https://6f1e-123-58-221-57.ngrok.io/oauth_sdk/app_uri', the app uri endpoint will be /oauth_sdk/app_uri
	r.GET("/oauth_sdk/app_uri", func(c *gin.Context) {
		params := getParams(c)
		var opts []co.AuthCodeOption
		c.Redirect(302, oauth.AuthCodeURL(params.Get("shop"), opts...))
	})

	// redirect uri endpoint
	// if your redirect uri is 'https://6f1e-123-58-221-57.ngrok.io/oauth_sdk/redirect_uri', the redirect uri endpoint will be /oauth_sdk/redirect_uri
	r.GET("/oauth_sdk/redirect_uri", func(c *gin.Context) {
		params := getParams(c)
		token, err := oauth.Exchange(context.Background(), params.Get("shop"), params.Get("code"))
		if err != nil {
			c.String(500, err.Error())
			c.Abort()
			return
		}
		c.JSON(200, token)
	})

	r.Run(":8080") // :your app server port
}

func getParams(c *gin.Context) url.Values {
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
