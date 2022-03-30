<!-- vscode-markdown-toc -->
* 1. [Introduce](#Introduce)
* 2. [Quick start](#Quickstart)
	* 2.1. [Create app](#Createapp)
	* 2.2. [Configure the app](#Configuretheapp)
	* 2.3. [Start SDK DemoApp](#sdk-demo-app)
	* 2.4. [Validation](#check)
* 3. [The directory structure](#dir)
* 4. [About the use of functional functions](#function)
	* 4.1. [AuthCodeURL](#AuthCodeURL)
	* 4.2. [ValidShop&&SignatureValid](#ValidShop&&SignatureValid)
	* 4.3. [Exchange](#Exchange)
	* 4.4. [RefreshToken](#RefreshToken)
* 5. [About the use of middleware](#middleware)
    * 5.1. [SetRequestPath](#SetRequestPath)
    * 5.2. [SetCallbackPath](#SetCallbackPath)
    * 5.3. [IgnoreState](#IgnoreState)
    * 5.4. [Middleware version of demo app](#MiddlewareDemoApp)

<!-- vscode-markdown-toc-config
	numbering=true
	autoSave=true
	/vscode-markdown-toc-config -->
<!-- /vscode-markdown-toc -->

<!-- [toc] -->
[中文版](./README.md)
##  1. <a name='Introduce'></a>Introduce

This project is a Go SDK for Shoplazza developers to complete authentication without having to understand too much Oauth2 process.

Read the documentation for the Shoplazza certification process Standard OAuth process [Standard OAuth process](https://helpcenter.shoplazza.com/hc/zh-cn/articles/4408686586137#h_01FM4XX2CX746V3277HB7SPGTN)


##  2. <a name='Quickstart'></a>Quick start
###  2.1. <a name='Createapp'></a>Create app
Read the documentation for creating the app [Building Public App](https://helpcenter.shoplazza.com/hc/zh-cn/articles/4409360434201)

###  2.2. <a name='Configuretheapp'></a>Configure the app
Read the documentation for configuring the app [Manage Your App](https://helpcenter.shoplazza.com/hc/zh-cn/articles/4409476265241)

###  2.3. <a name='sdk-demo-app'></a>Start SDK DemoApp
1、Import Oauth SDK Go
```curl
go get -u https://github.com/Shoplazza/oauth-sdk-go
```

2、Enter the example/basic directory
```
cd /example/basic
```

3、Fill in the values of ClientID, ClientSecret, RedirectUri and Scopes of your app into the config structure
```go
import (
  co "github.com/Shoplazza/oauth-sdk-go"
)

oauth := &co.Config{
  ClientID:     "WVpHNkENUL9CBbDjGO_po9tfG02XW2Z-X54M4LObfDs",     // your app client id
  ClientSecret: "DdJNhsopKxAWDHjqjI1rpQZW17Fp6GXrHhC0IgwdXag", // your app client secret
  Endpoint:     shoplazza.Endpoint,
  RedirectURI:  "https://6f1e-123-58-221-57.ngrok.io/oauth_sdk/redirect_uri",      // your app redirect uri, you need to replace the 6f1e-123-58-221-57.ngrok.io with your service domain, and ensure that the domain name is externally accessible
  Scopes:       []string{"read_shop", "write_shop"}, // the permissions you want
}
```
4、Fill in the correct App Uri path and Redirect Uri path
```go
// app uri path
// if your redirect uri is 'https://6f1e-123-58-221-57.ngrok.io/oauth_sdk/app_uri', the app uri path will be '/oauth_sdk/app_uri'
r.GET("/oauth_sdk/app_uri", func(c *gin.Context) {
    params := getParams(c)
    var opts []co.AuthCodeOption
    c.Redirect(302, oauth.AuthCodeURL(params.Get("shop"), opts...))
})
```
```go
// redirect uri path
// if your redirect uri is 'https://6f1e-123-58-221-57.ngrok.io/oauth_sdk/redirect_uri', the redirect uri path will be '/oauth_sdk/redirect_uri'
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
5、Set the port of the service
```go
r.Run(":8080") // here is your service port. 8080 is just an example
```

###  2.4. <a name='check'></a>Validation
Now that all configurations of demo app have been completed, start app for verification

1、Execute in the baisc directory
```curl
go run demo_app.go
```

2、Ensure that the installation process is smooth
Please read this step first [Testing public app](https://helpcenter.shoplazza.com/hc/en-us/articles/4409360434201-Building-Public-App#h_01FM7GXEAM5VPXTK6PJSA9MFWC)

Installation process:

If you have development store, go to  [Partner Center](https://partners.shoplazza.com/)->Apps->Apps List->Manage Apps->Test Apps  with the above Shoplazza account, select store to install the app and jump to the authorized installation page of the store for testing.

The page responds to the normal authorization result:
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

If you encounter the following page, you need to replace the 6f1e-123-58-221-57.ngrok.io with your service domain!!!
![报错图片](https://cdn.shoplazza.com/59c635bd81b66755ab9a64d698af0900.jpg)

##  3. <a name='dir'></a>The directory structure
```shell
.
├── example
│   ├── basic   
│   │     └── demo-app.go         // Demo app, Fill in the configuration to start the service
│   ├── middleware    
│   │     └── demo-app.go         // Demo app of middleware version,Fill in the configuration to start the service
├── internal
│   └── token.go           
├── shoplazza
│   └── shoplazza.go
├── go.mod
├── go.sum
├── middleware.go           // Provide gin middleware,there are already implemented app URI and redirect URI methods,help you quickly complete the authorization process
├── oauth2.go               // It includes many methods such as AuthCodeURL,ValidShop,Exchange,help you quickly complete the authorization process
├── REAEDME-ch.md
├── README.md
├── token.go    
```

##  4. <a name='function'></a>About the use of functional functions

###  4.1. <a name='AuthCodeURL'></a>AuthCodeURL 
To show the App's prompt page for the merchant to start with,Shoplazza will first call the `app uri path` provided by your app service,you need to redirect the parameters defined below to the following URL in this endpoint:
```curl
https://{store_name}.myshoplaza.com/admin/oauth/authorize?client_id={client_id}&scope={scopes}&redirect_uri={redirect_uri}&response_type={response_type}&state={state}
```
- store_name: The name of merchant's store.
- client_id: App client id
- scopes: A space separated list of scopes. For example, to write orders and read customers, use scope="write_order read_customer".
- redirect_uri: The URL to which a merchant is redirected after authorizing the app.
- response_type: The response type of OAuth 2.0 process, here we need to fill in "code"
- state: A random value, use to prevent CSRF attacks.

You can use Go OAuth SDK to quickly assemble this URL, see example below:
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
When the merchant clicks the install button in the prompt, they will be redirected to the `redirect uri path` of your app service,see example below:
```curl
http://example.com/some/redirect_uri?code={authorization_code}&shop={store_name}.myshoplaza.com&hmac={hmac}
```

Before we continue, make sure your app performs the following security checks. If any of the checks faill, then your app must reject the request with an error, and must not continue.
- The `hmac` is valid and signed by Shoplazza
- The `shop` parameter is a valid shop hostname, ends with `myshoplaza.com`

For Security Checks, Go OAuth SDK also has corresponding methods, you can quickly verify hmac and shop parameter by SDK, see example below:

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
If all security checks pass, then you can exchange the authorization_code for a permanent access token by sending a request to the shop's access_token endpoint:

```
POST https://{store_name}.myshoplaza.com/admin/oauth/token
```
In this request, store_name is the name of the merchant's store and alongs with the following parameters:
- client_id: app client id
- client_secret: The Client secret key for the app.
- code: The authorization_code provided in the redirect.
- grant_type: The grant type of OAuth 2.0 process, please fill in "authorization_code" here.
- redirect_uri: The redirect_uri of the app.

The server responds with an access token:
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
- token_type: It will just return Bearer.
- expires_at: access_token expired time, in timestamp.
- refresh_token: The refresh token used to refresh the access_token if needed.
- access_token: The correct access_token.
- store_id: Store's ID in Shoplazza
- store_name: Store name

Similarly, you can quickly get an access token by Go OAuth SDK, see example below:
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
After access_token expired, The app need to call endpoint to retrieve a new access_token and a new refresh_token ( Please save it into your app and you are gonna need it later)

Similarly, you can quickly refresh the access token by Go OAuth SDK, see example below:
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

##  5. <a name='middleware'></a>About the use of middleware
Gin middleware will intercept the requests of `/auth/shoplazza` and `/auth/shoplazza/callback` by default:
- `/auth/shoplazza?shop=xx.myshoplaza.com` : Asked this URL will be redirected to the https://xx.myshoplaza.com/admin/oauth/authorize to initiate the authorization process
- `/auth/shoplazza/callback` : Intercepts authorization callback requests and automatically replaces tokens with codes in callback requests

###  5.1. <a name='SetRequestPath'></a>SetRequestPath
This method has encapsulated verification of store domain name, state and redirect to https://xx.myshoplaza.com/admin/oauth/authorize, you only need to set the `app uri path`

See example below:
```go
import (
  co "github.com/Shoplazza/oauth-sdk-go"
)

r := gin.New()
oauth := &co.Config{ xxx... }
oauthMid := co.NewGinMiddleware(oauth)
oauthMid.SetRequestPath("app uri path") // custom request path
r.Use(oauthMid.Handler())
```

###  5.2. <a name='SetCallbackPath'></a>SetCallbackPath
This method has encapsulated the verification of store domain name, state and HMAC, and exchange token, and saves the obtained token to the context. You only need to set `redirect uri path`, implement `redirect URI path` logic, and then obtain the token directly from the context.

See example below:
```go
import (
  co "github.com/Shoplazza/oauth-sdk-go"
)

r := gin.New()
oauth := &co.Config{ xxx... }
oauthMid := co.NewGinMiddleware(oauth)
oauthMid.SetCallbackPath("redirect uri path") // custom callback path
r.Use(oauthMid.Handler())
// implement your redirect uri logic,obtain the token directly from the context
r.GET("redirect uri path", func(c *gin.Context) {
    t, _ := c.Get("oauth2.token")
    token := t.(oauth2.Token)
    log.Println("example = ", token)
})
```

###  5.3. <a name='IgnoreState'></a>IgnoreState
This method enables you to ignore the verification of state in the implementation logic of `redirect uri path`

See example below:
```go
import (
  co "github.com/Shoplazza/oauth-sdk-go"
)

r := gin.New()
oauth := &co.Config{ xxx... }
oauthMid := co.NewGinMiddleware(oauth)
oauthMid.IgnoreState() // ignore state
r.Use(oauthMid.Handler())
```
###  5.4. <a name='MiddlewareDemoApp'></a>Middleware version of demo app
If you want to try the demo app of the middleware version, we have prepared it for you. You only need to configure it as you configure the example demo app to start it quickly
See example below:

1、Enter the example/middleware directory
```
cd /example/middleware
```

2、Fill in the values of ClientID, ClientSecret, RedirectUri and Scopes of your app into the config structure
```go
import (
  co "github.com/Shoplazza/oauth-sdk-go"
)

oauth := &co.Config{
  ClientID:     "WVpHNkENUL9CBbDjGO_po9tfG02XW2Z-X54M4LObfDs",     // your app client id
  ClientSecret: "DdJNhsopKxAWDHjqjI1rpQZW17Fp6GXrHhC0IgwdXag", // your app client secret
  Endpoint:     shoplazza.Endpoint,
  RedirectURI:  "https://6f1e-123-58-221-57.ngrok.io/oauth_sdk/redirect_uri",      // your app redirect uri, you need to replace the 6f1e-123-58-221-57.ngrok.io with your service domain, and ensure that the domain name is externally accessible
  Scopes:       []string{"read_shop", "write_shop"}, // the permissions you want
}
```

3、Fill in the correct App Uri path and Redirect Uri path
```go
oauthMid.SetRequestPath("/oauth_sdk/app_uri")
oauthMid.SetCallbackPath("/oauth_sdk/redirect_uri")
```

4、Fill in and implement the redirect URI path, and do the processing after obtaining the token, such as storing it in DB
```go
r.GET("/oauth_sdk/redirect_uri", func(c *gin.Context) {
  t, _ := c.Get("oauth2.token")
  token := t.(*co.Token)
  c.JSON(200, token)
})
```

5、Set the port of the service
```go
r.Run(":8080") // here is your service port. 8080 is just an example
```

6、Start
```
go run demo_app.go
```