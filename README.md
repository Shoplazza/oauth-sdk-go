# OAuth2 for Go

oauth2 package contains a client implementation for OAuth 2.0 spec.

## Installation
...

## Usage

```go
oauth := &oauth2.Config{
    ClientID:     "xxx",
    ClientSecret: "xxx",
    RedirectURI: "https://app.com/auth/shoplazza/callback",
    Scopes:      []string{"read_product", "write_product"},
    Endpoint:    shoplazza.Endpoint,
    Domain:      "preview.shoplazza.com", // If not set, production domain name will be used by default: myshoplaza.com
}

// Use code to exchange tokens when authorizing callbacks
token, err := oauth.Exchange(context.Background(), "xxx.myshoplaza.com", "code")

// Refresh token when token expires
token, err := oauth.RefreshToken(context.Background(), "xxx.myshoplaza.com", "refresh token")
```

#### In Gin
Gin middleware will intercept requests from `/auth/shoplazza` and `/auth/shoplazza/callback` URLs by default:
- `/auth/shoplazza?shop=xx.myshoplaza.com` : When requesting this URL, it will redirect to https://xx.myshoplaza.com/admin/oauth/authorize to initiate the authorization process
- `/auth/shoplazza/callback` : Intercept the authorization callback request and automatically replace the token with the code in the callback request

```go
r := gin.New()
oauth := &oauth2.Config{ xxx... }

oauthMid := oauth2.NewGinMiddleware(oauth)
oauthMid.SetCallbackPath("xxx") // custom callback path
oauthMid.SetCallbackFunc(func(c *gin.Context) { // custom callback function
    // ....
})
oauthMid.SetRequestPath("xxx") // custom request path
oauthMid.SetRequestFunc(func(c *gin.Context) { // Custom request function
    // ....
})

r.Use(oauthMid.Handler())

r.GET(oauth2.DefaultCallbackPath, func(c *gin.Context) {
    t, _ := c.Get("oauth2.token")
    token := t.(oauth2.Token)
    log.Println("example = ", token)
})
_ = r.Run(":8080")
```
#### Attach a simple Demo APP
```go
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	co "github.com/shoplazza-os/oauth-sdk-go"
	"github.com/shoplazza-os/oauth-sdk-go/shoplazza"
)

type Token struct {
	Id          int       `json:"id"`
	AccessToken string    `json:"access_token"`
	Shop        string    `json:"shop"`
	StoreId     string    `json:"store_id"`
	ExpiredAt   time.Time `json:"expired_at"`
}

type ShopResult struct {
	Shop struct {
		Id            string `json:"id"`
		Name          string `json:"name"`
		CountryCode   string `json:"country_code"`
		ProvinceCode  string `json:"province_code"`
		Address1      string `json:"address1"`
		Address2      string `json:"address2"`
		Phone         string `json:"phone"`
		PrimaryLocale string `json:"primary_locale"`
		Currency      string `json:"currency"`
	} `json:"shop"`
}

func main() {
	oauth := &co.Config{
		ClientID:     "s1Ip1WxpoEAHtPPzGiP2rK2Az-P07Nie7V97hRKigl4",
		ClientSecret: "0LFJcNqVb2Z1nVt9xT72vOo0sTWd6j8wVX60Y5xdzZ0",
		Endpoint:     shoplazza.Endpoint,
		RedirectURI:  "https://3830-43-230-206-233.ngrok.io/oauth_sdk/redirect_uri/",
		Scopes:       []string{"read_shop"},
	}

	db, err := gorm.Open("mysql", "root:123456@/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	db.LogMode(true)
	defer db.Close()

	oauthMid := co.NewGinMiddleware(oauth)
	// Set whether to do state verification
	oauthMid.IgnoreState()
	// You can use the encapsulated function, or you can use a custom callback function
	oauthMid.SetCallbackPath("/oauth_sdk/redirect_uri/")                                
	// You can use the encapsulated function, or you can use a custom request function
	oauthMid.SetRequestPath("/oauth_sdk/app_uri")                                         
    // This function is to do something after getting the token. If you use a custom callback function, this does not need to be implemented
	oauthMid.SetAccessTokenHandlerFunc(func(c *gin.Context, shop string, token *co.Token) { 
		// store-id is stored in the cookie
		c.SetCookie("store_id", token.StoreInfo.ID, 3600, "/", "https://3830-43-230-206-233.ngrok.io", false, true)
		// store-id, shop, access-token and other information are stored in DB
		t := Token{
			AccessToken: token.AccessToken,
			Shop:        shop,
			StoreId:     token.StoreInfo.ID,
			ExpiredAt:   token.ExpiresAt,
		}
		db.Save(&t)
		c.String(200, "save access-token success")
	})

	r := gin.Default()
	r.Use(oauthMid.Handler())
	// [Important] If oauthMid.SetAccessTokenHandlerFunc is used, the RedirectTrailingSlash of gin needs tto be set to false
	r.RedirectTrailingSlash = false

	r.GET("/open_api/test", func(c *gin.Context) {
		var req *http.Request
		httpclient := &http.Client{Transport: &http.Transport{
			TLSHandshakeTimeout: 3 * time.Second,
		},
			Timeout: 6 * time.Second,
		}

		storeId, _ := c.Cookie("store_id")
		var t Token
		if err := db.Table("tokens").Where("store_id = ?", storeId).Find(&t).Error; err != nil {
			c.JSON(404, "access token is not founded")
			c.Abort()
			return
		}

		req, _ = http.NewRequest("GET", fmt.Sprintf("https://%s/openapi/2020-07/shop", t.Shop), nil)
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Access-Token", t.AccessToken)

		resp, err := httpclient.Do(req)
		if err != nil {
			c.JSON(resp.StatusCode, err)
			c.Abort()
			return
		}

		defer resp.Body.Close()
		var shop ShopResult
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(resp.StatusCode, err)
			c.Abort()
			return
		}

		err = json.Unmarshal(content, &shop)
		if err != nil {
			log.Println(err)
		}

		c.JSON(resp.StatusCode, shop)
	})

	r.Run(":8080")
}
```