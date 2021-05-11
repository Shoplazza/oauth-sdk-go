package shoplazza

import (
	"oauth2"
)

var Endpoint = oauth2.Endpoint{
	AuthURL:  "/callback/shoplazza/oauth",
	TokenURL: "/admin/oauth/token",
}
