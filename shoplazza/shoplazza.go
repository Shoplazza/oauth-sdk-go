package shoplazza

import (
	oauth2 "gitlab.shoplazza.site/common/common-oauth2"
)

var Endpoint = oauth2.Endpoint{
	AuthURL:  "/admin/oauth/authorize",
	TokenURL: "/admin/oauth/token",
}
