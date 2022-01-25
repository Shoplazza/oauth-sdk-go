package shoplazza

import (
	oauth2 "oauth-sdk-go"
)

var Endpoint = oauth2.Endpoint{
	AuthURL:  "/admin/oauth/authorize",
	TokenURL: "/admin/oauth/token",
}
