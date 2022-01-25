package shoplazza

import (
	oauth2 "github.com/shoplazza-os/oauth-sdk-go"
)

var Endpoint = oauth2.Endpoint{
	AuthURL:  "/admin/oauth/authorize",
	TokenURL: "/admin/oauth/token",
}
