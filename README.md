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