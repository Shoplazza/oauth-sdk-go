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
    Domain:      "preview.shoplazza.com", // 不设置的话默认使用美服域名：myshoplaza.com
}

// 授权回调时使用 code 交换 token
token, err := oauth.Exchange(context.Background(), "xxx.myshoplaza.com", "code")

// token 过期时刷新 token
token, err := oauth.RefreshToken(context.Background(), "xxx.myshoplaza.com", "refresh token")
```

#### In Gin
Gin middleware 会默认拦截 `/auth/shoplazza` 以及 `/auth/shoplazza/callback` 两个 URL 的请求:
- `/auth/shoplazza?shop=xx.myshoplaza.com` : 请求此 URL 时，会重定向到 https://xx.myshoplaza.com/admin/oauth/authorize 去发起授权流程
- `/auth/shoplazza/callback` : 拦截授权回调请求，自动将回调请求中的 code 替换 token

```go
r := gin.New()
oauth := &oauth2.Config{ xxx... }

oauthMid := oauth2.NewGinMiddleware(oauth)
oauthMid.SetCallbackPath("xxx") // 自定义 callback path
oauthMid.SetCallbackFunc(func(c *gin.Context) { // 自定义 callback 处理函数
    // ....
})
oauthMid.SetRequestPath("xxx") // 自定义 request path
oauthMid.SetRequestFunc(func(c *gin.Context) { // 自定义 request 处理函数
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
#### 附上一个简易的 Demo APP 代码
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
	co "gitlab.shoplazza.site/common/common-oauth2"
	"gitlab.shoplazza.site/common/common-oauth2/shoplazza"
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
	// 设置是否做 state 校验
	oauthMid.IgnoreState()
	// 可以使用封装好的方法，也可以自定义 callback 处理函数
	oauthMid.SetCallbackPath("/oauth_sdk/redirect_uri/")                                
	// 可以使用封装好的方法，也可以自定义 request 处理函数
	oauthMid.SetRequestPath("/oauth_sdk/app_uri")                                         
    // 该方法是拿到token后的后续处理逻辑，如果使用自定义 callback 处理函数，这个就不需要实现了
	oauthMid.SetAccessTokenHandlerFunc(func(c *gin.Context, shop string, token *co.Token) { 
		// store-id 存 cookie
		c.SetCookie("store_id", token.StoreInfo.ID, 3600, "/", "https://3830-43-230-206-233.ngrok.io", false, true)
		// store-id、shop、access-token等信息存DB
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
	// 【重要】如果使用了 oauthMid.SetAccessTokenHandlerFunc，需要关闭 gin 的自动重定向功能
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