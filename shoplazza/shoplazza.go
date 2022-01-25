package shoplazza

import (
	oauth2 "github.com/shoplazza-os/Oauth-SDK-Go"
)

var Endpoint = oauth2.Endpoint{
	AuthURL:  "/admin/oauth/authorize",
	TokenURL: "/admin/oauth/token",
}
