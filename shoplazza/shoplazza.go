package shoplazza

import (
	oauth2 "common_oauth2"
)

var Endpoint = oauth2.Endpoint{
	AuthURL:  "/admin/oauth/authorize",
	TokenURL: "/admin/oauth/token",
}
