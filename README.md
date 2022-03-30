<!-- vscode-markdown-toc -->
* 1. [介绍](#introduce)
* 2. [快速启动](#-1)
	* 2.1. [创建App](#app)
	* 2.2. [配置App](#app-1)
	* 2.3. [启动SDK DemoApp](#sdk-demo-app)
	* 2.4. [验证](#check)
* 3. [目录结构](#dir)
* 4. [关于功能函数的使用](#function)
	* 4.1. [AuthCodeURL](#AuthCodeURL)
	* 4.2. [ValidShop&&SignatureValid](#ValidShop&&SignatureValid)
	* 4.3. [Exchange](#Exchange)
	* 4.4. [RefreshToken](#RefreshToken)
* 5. [中间件的使用](#middleware)
    * 5.1. [SetRequestPath](#SetRequestPath)
    * 5.2. [SetCallbackPath](#SetCallbackPath)
    * 5.3. [IgnoreState](#IgnoreState)
    * 5.4. [中间件版本的DemoApp](#MiddlewareDemoApp)

<!-- vscode-markdown-toc-config
	numbering=true
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->

<!-- [toc] -->
[English version](./README-en.md)
##  1. <a name='introduce'></a>介绍

本项目是为了shoplazza的开发者可以不需要理解过多的Oauth2流程完成认证操作的一款Go语言开发的sdk。

关于shoplazza认证流程请阅读文档 [标准的OAuth流程](https://helpcenter.shoplazza.com/hc/zh-cn/articles/4408686586137#h_01FM4XX2CX746V3277HB7SPGTN)


##  2. <a name='-1'></a>快速启动
###  2.1. <a name='app'></a>创建app
关于创建app 请阅读文档 [构建公用App](https://helpcenter.shoplazza.com/hc/zh-cn/articles/4409360434201)

###  2.2. <a name='app-1'></a>配置app
关于配置app 请阅读文档 [管理你的App](https://helpcenter.shoplazza.com/hc/zh-cn/articles/4409476265241)

###  2.3. <a name='sdk-demo-app'></a>启动SDK DemoApp
1、引入Oauth SDK Go
```curl
go get -u https://github.com/Shoplazza/oauth-sdk-go
```

2、进入example/basic目录
```
cd /example/basic
```

3、填入你的App的ClientID、ClientSecret、RedirectUri、Scopes的值到Config结构体中
```go
import (
  co "github.com/Shoplazza/oauth-sdk-go"
)

oauth := &co.Config{
  ClientID:     "WVpHNkENUL9CBbDjGO_po9tfG02XW2Z-X54M4LObfDs", // App的ClientId
  ClientSecret: "DdJNhsopKxAWDHjqjI1rpQZW17Fp6GXrHhC0IgwdXag", // App的ClientSecret
  Endpoint:     shoplazza.Endpoint,
  RedirectURI:  "https://6f1e-123-58-221-57.ngrok.io/oauth_sdk/redirect_uri", // App的RedirectUri, 你需要把域名 6f1e-123-58-221-57.ngrok.io 换成自己服务的域名，并且保证域名是对外可访问的
  Scopes:       []string{"read_shop", "write_shop"},                          // 想要申请的权限
}
```
4、填入正确的App Uri path和Redirect Uri path
```go
// app uri path
// 如果你完整的AppUri是 'https://6f1e-123-58-221-57.ngrok.io/oauth_sdk/app_uri', 那app uri path就是 '/oauth_sdk/app_uri'
r.GET("/oauth_sdk/app_uri", func(c *gin.Context) {
    params := getParams(c)
    var opts []co.AuthCodeOption
    c.Redirect(302, oauth.AuthCodeURL(params.Get("shop"), opts...))
})
```
```go
// redirect uri path
// 如果你完整的RedirectUri是 'https://6f1e-123-58-221-57.ngrok.io/oauth_sdk/redirect_uri', 那app uri path就是 '/oauth_sdk/redirect_uri'
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
```
5、设置服务的端口
```go
r.Run(":8080") // 这里配置成你服务的端口，8080只是示例
```

###  2.4. <a name='check'></a>验证
此时已经完成了Demo App所有的配置，那么接下来启动App进行验证

1、在basic目录下执行
```curl
go run demo_app.go
```

2、确定安装流程顺畅
这个步骤请先阅读 [测试公共App](https://helpcenter.shoplazza.com/hc/zh-cn/articles/4409360434201#h_01FM7BPX2QBPB9ZWQZM80GTH4C)

按照流程：

前往 [合作伙伴中心](https://partners.shoplazza.com/)->App->App列表->管理App->测试App 入口，选择该店铺安装App，即可跳转至该店铺的授权安装页面.

正常授权结果后页面会返回:
```json
{
    "AccessToken": "lzhOS5Gl3tfVDSQZ8pBVQuMyCD24PMp9sUGZXMtW3b4",
    "TokenType": "Bearer",
    "RefreshToken": "Um8UXkriF_5-naBX0cYw31c-jiCdKBi1v-hPnY7DMdg",
    "ExpiresAt": "2023-03-24T15:37:18+08:00",
    "StoreInfo": {
        "ID": "168705",
        "Locale": "zh-CN",
        "Name": "pjs"
    }
}
```

如果你遇到以下页面，你需要把域名 `6f1e-123-58-221-57.ngrok.io` 换成自己服务可访问的域名!!!
![报错图片](https://cdn.shoplazza.com/59c635bd81b66755ab9a64d698af0900.jpg)

##  3. <a name='dir'></a>目录结构
```shell
.
├── example
│   ├── basic   
│   │     └── demo-app.go         // demo app，填入配置即可启动服务
│   ├── middleware    
│   │     └── demo-app.go         // 中间件版本demo app，填入配置即可启动服务
├── internal
│   └── token.go           
├── shoplazza
│   └── shoplazza.go        // 包含发起授权的Endpoint和code换取token的Endpoint
├── go.mod
├── go.sum
├── middleware.go           // 提供 gin 中间件，里面已经有实现好的 app uri 和 redirect uri 方法，帮助你快速完成授权流程
├── oauth2.go               // 封装了许多参数校验、code换取token，构造AuthCodeURL等方法，帮助你快速完成授权流程
├── REAEDME-ch.md
├── README.md
├── token.go    
```

##  4. <a name='function'></a>关于功能函数的使用

###  4.1. <a name='AuthCodeURL'></a>AuthCodeURL 
为了向商家显示应用程序的提示页面，Shoplazza将首先调用你的App服务提供的`app uri path`，你需要在这个path里将以下定义的参数重定向到以下URL：
```curl
https://{store_name}.myshoplaza.com/admin/oauth/authorize?client_id={client_id}&scope={scopes}&redirect_uri={redirect_uri}&response_type={response_type}&state={state}
```
- store_name: 商户的店铺名称
- client_id: app client id
- scopes: 需要用到的权限列表，用法: "read_product read_customer"
- redirect_uri: 商户授权后重定向的URL地址
- response_type: Oauth2.0的返回类型，这里我们只需要填"code"
- state: 一个随机值，用来防止CSRF攻击

使用Oauth SDK，你能够快速组装这个URL，用法如下:
```go
import (
  co "github.com/Shoplazza/oauth-sdk-go"
  "github.com/Shoplazza/oauth-sdk-go/shoplazza"
)

oauth := &co.Config{
    ClientID:     "s1Ip1WxpoEAHtPPzGiP2rK2Az-P07Nie7V97hRKigl4",
    ClientSecret: "0LFJcNqVb2Z1nVt9xT72vOo0sTWd6j8wVX60Y5xdzZZ",
    Endpoint:     shoplazza.Endpoint,
    RedirectURI:  "https://3830-43-230-206-233.ngrok.io/oauth_sdk/redirect_uri/",
    Scopes:       []string{"read_shop"},
}
var opts []AuthCodeOption
oauth.AuthCodeURL("xxx.myshoplaza.com", opts...)
```

###  4.2. <a name='ValidShop&&SignatureValid'></a>ValidShop&&SignatureValid
当商家点击提示中的安装App按钮时，他们会被重定向到你的App服务的`redirect uri path`，例子如下：
```curl
http://example.com/some/redirect_uri?code={authorization_code}&shop={store_name}.myshoplaza.com&hmac={hmac}
```

你需要在这个path中执行以下安全检查，如果任何检查失败，那么你的应用程序必须以错误拒绝请求
- `hmac`必须有效且由Shoplazza分配的ClientSecret签名
- `shop`参数必须是一个有效的店铺域名，以`myshoplaza.com`结尾

使用Oauth SDK，你能够快速完成安全检查，用法如下:
```go
import (
  co "github.com/Shoplazza/oauth-sdk-go"
  "github.com/Shoplazza/oauth-sdk-go/shoplazza"
)

oauth := &co.Config{
    ClientID:     "s1Ip1WxpoEAHtPPzGiP2rK2Az-P07Nie7V97hRKigl4",
    ClientSecret: "0LFJcNqVb2Z1nVt9xT72vOo0sTWd6j8wVX60Y5xdzZZ",
    Endpoint:     shoplazza.Endpoint,
    RedirectURI:  "https://3830-43-230-206-233.ngrok.io/oauth_sdk/redirect_uri/",
    Scopes:       []string{"read_shop"},
}


var redirectUrl = "http://example.com/some/redirect_uri?code={authorization_code}&shop={store_name}.myshoplaza.com&hmac={hmac}"
query := strings.Split(redirectUrl, "?")
params, _ := url.ParseQuery(query[1])
oauth.ValidShop(params.Get("shop"))        // verify shop parameter
oauth.SignatureValid(params)               // verify hmac
```

###  4.3. <a name='Exchange'></a>Exchange
如果所有安全检查都通过，那么您可以通过向店铺的获取token endpoint发送请求，用授权码code来换取token：
```
POST https://{store_name}.myshoplaza.com/admin/oauth/token
```
在这个请求中，store_name是商家商店的名称，并带有以下参数：
- client_id: app client id
- client_secret: app client 密钥
- code: redirect uri 参数中的授权码
- grant_type: Oauth2.0的授权类型，这里我们只需要填"authorization_code"
- redirect_uri: app redirect uri

这个接口会返回以下信息:
```json
{
  "token_type": "Bearer",
  "expires_at": 1550546245,
  "access_token": "eyJ0eXAiOiJKV1QiLCJh",
  "refresh_token": "def502003d28ba08a964e",
  "store_id": "2",
  "store_name": "xiong1889"
}
```
- token_type: 这里一般会返回 "Bearer"
- expires_at: access token的过期时间，返回时间戳
- refresh_token: token过期后，需要使用refresh token来刷新access token
- access_token: 正确的access token
- store_id: 店铺ID
- store_name: 店铺名称

使用Oauth SDK，你能够快速获取access token，用法如下:
```go
import (
  co "github.com/Shoplazza/oauth-sdk-go"
  "github.com/Shoplazza/oauth-sdk-go/shoplazza"
)

oauth := &co.Config{
    ClientID:     "s1Ip1WxpoEAHtPPzGiP2rK2Az-P07Nie7V97hRKigl4",
    ClientSecret: "0LFJcNqVb2Z1nVt9xT72vOo0sTWd6j8wVX60Y5xdzZZ",
    Endpoint:     shoplazza.Endpoint,
    RedirectURI:  "https://3830-43-230-206-233.ngrok.io/oauth_sdk/redirect_uri/",
    Scopes:       []string{"read_shop"},
}
token, err := oauth.Exchange(context.Background(),"xxx.myshoplaza.com", "code"))
```

###  4.4. <a name='RefreshToken'></a>RefreshToken
token过期后，使用refresh token来刷新access token
使用Oauth SDK，你能够快速刷新access token，用法如下:
```go
import (
  co "github.com/Shoplazza/oauth-sdk-go"
  "github.com/Shoplazza/oauth-sdk-go/shoplazza"
)

oauth := &co.Config{
    ClientID:     "s1Ip1WxpoEAHtPPzGiP2rK2Az-P07Nie7V97hRKigl4",
    ClientSecret: "0LFJcNqVb2Z1nVt9xT72vOo0sTWd6j8wVX60Y5xdzZZ",
    Endpoint:     shoplazza.Endpoint,
    RedirectURI:  "https://3830-43-230-206-233.ngrok.io/oauth_sdk/redirect_uri/",
    Scopes:       []string{"read_shop"},
}
token, err := oauth.RefreshToken(context.Background(), "xxx.myshoplaza.com", "refresh token")
```


##  5. <a name='middleware'></a>中间件的使用
Gin middleware 会默认拦截 `/auth/shoplazza` 以及 `/auth/shoplazza/callback` 两个 URL 的请求:
- `/auth/shoplazza?shop=xx.myshoplaza.com` : 请求此 URL 时，会重定向到 https://xx.myshoplaza.com/admin/oauth/authorize 去发起授权流程
- `/auth/shoplazza/callback` : 拦截授权回调请求，自动将回调请求中的 code 替换 token

###  5.1. <a name='SetRequestPath'></a>SetRequestPath

该方法已经封装了校验店铺域名，校验state，重定向到 https://xx.myshoplaza.com/admin/oauth/authorize 的功能，你只需要设置`app uri path`

使用方法：
```go
import (
  co "github.com/Shoplazza/oauth-sdk-go"
)

r := gin.New()
oauth := &co.Config{ xxx... }
oauthMid := co.NewGinMiddleware(oauth)
oauthMid.SetRequestPath("app uri path") // 自定义 request path
r.Use(oauthMid.Handler())
```

###  5.2. <a name='SetCallbackPath'></a>SetCallbackPath

该方法已经封装了校验店铺域名，校验state，校验hmac，并调用code换取token的接口，并将获取的token存到context中，你只需要设置`redirect uri path`，并实现`redirect uri path`，然后直接从context获取token

使用方法：
```go
import (
  co "github.com/Shoplazza/oauth-sdk-go"
)

r := gin.New()
oauth := &co.Config{ xxx... }
oauthMid := co.NewGinMiddleware(oauth)
oauthMid.SetCallbackPath("redirect uri path") // 自定义 callback path
r.Use(oauthMid.Handler())
// 实现你的redirect uri逻辑，从context直接获取token
r.GET("redirect uri path", func(c *gin.Context) {
    t, _ := c.Get("oauth2.token")
    token := t.(oauth2.Token)
    log.Println("example = ", token)
})
```

###  5.3. <a name='IgnoreState'></a>IgnoreState

该方法可以使你在`redirect uri path`的实现逻辑中，忽略state的校验

使用方法：
```go
import (
  co "github.com/Shoplazza/oauth-sdk-go"
)

r := gin.New()
oauth := &co.Config{ xxx... }
oauthMid := co.NewGinMiddleware(oauth)
oauthMid.IgnoreState() // 设置是否做 state 校验
r.Use(oauthMid.Handler())
```
###  5.4. <a name='MiddlewareDemoApp'></a>中间件版本的DemoApp

如果你想尝试中间件版本的Demo App，我们已经为你准备好了，你只需要和配置启动Example Demo App一样对它进行配置，即可快速启动

使用方法：

1、进入中间件Demo App目录
```
cd /example/middleware
```

2、填入你的App的ClientID、ClientSecret、RedirectUri、Scopes的值到Config结构体中
```go
import (
  co "github.com/Shoplazza/oauth-sdk-go"
)

oauth := &co.Config{
  ClientID:     "WVpHNkENUL9CBbDjGO_po9tfG02XW2Z-X54M4LObfDs", // App的ClientId
  ClientSecret: "DdJNhsopKxAWDHjqjI1rpQZW17Fp6GXrHhC0IgwdXag", // App的ClientSecret
  Endpoint:     shoplazza.Endpoint,
  RedirectURI:  "https://6f1e-123-58-221-57.ngrok.io/oauth_sdk/redirect_uri", // App的RedirectUri, 你需要把域名 6f1e-123-58-221-57.ngrok.io 换成自己服务可访问的域名
  Scopes:       []string{"read_shop", "write_shop"},                          // 想要申请的权限
}
```

3、填入正确的App Uri path和Redirect Uri path
```go
oauthMid.SetRequestPath("/oauth_sdk/app_uri")
oauthMid.SetCallbackPath("/oauth_sdk/redirect_uri")
```

4、填入并实现Redirect Uri path，做获取token后的处理，如：存入DB
```go
r.GET("/oauth_sdk/redirect_uri", func(c *gin.Context) {
  t, _ := c.Get("oauth2.token")
  token := t.(*co.Token)
  c.JSON(200, token)
})
```

5、设置服务的端口
```go
r.Run(":8080") // 这里配置成你服务的端口，8080只是示例s
```

6、启动
```
go run demo_app_middleware.go
```