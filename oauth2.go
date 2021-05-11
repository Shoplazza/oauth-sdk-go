package oauth2

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"oauth2/internal"
	"regexp"
	"strings"
)

var DefaultDomain = "myshoplaza.com"

type Config struct {
	// application's ID.
	ClientID string

	// application's secret.
	ClientSecret string

	Endpoint Endpoint

	RedirectURI string

	// Scope specifies optional requested permissions.
	Scopes []string
	Domain string
}

type Endpoint struct {
	AuthURL  string
	TokenURL string
}

// request_path    /auth/shoplazza?shop=huangpuping.shoplazza.com
// callback_path   /auth/callback/shoplazza?shop=huangpuping.shoplazza.com

// An AuthCodeOption is passed to Config.AuthCodeURL.
type AuthCodeOption interface {
	setValue(url.Values)
}

type setParam struct{ k, v string }

func (p setParam) setValue(m url.Values) { m.Set(p.k, p.v) }

// SetAuthURLParam builds an AuthCodeOption which passes key/value parameters to a provider's authorization endpoint.
func SetAuthURLParam(key, value string) AuthCodeOption {
	return setParam{key, value}
}

func (c *Config) AuthCodeURL(shop string, opts ...AuthCodeOption) string {
	authUrl := fmt.Sprintf("%s%s?", c.fixSite(shop), c.Endpoint.AuthURL)

	var buf bytes.Buffer
	buf.WriteString(authUrl)

	v := url.Values{
		"response_type": {"code"},
		"client_id":     {c.ClientID},
	}
	if c.RedirectURI != "" {
		v.Set("redirect_uri", c.RedirectURI)
	}
	if len(c.Scopes) > 0 {
		v.Set("scope", strings.Join(c.Scopes, " "))
	}
	for _, opt := range opts {
		opt.setValue(v)
	}

	buf.WriteString(v.Encode())
	return buf.String()
}

func (c *Config) Exchange(ctx context.Context, shop, code string, opts ...AuthCodeOption) (*Token, error) {
	v := url.Values{
		"grant_type": {"authorization_code"},
		"code":       {code},
	}
	if c.RedirectURI != "" {
		v.Set("redirect_uri", c.RedirectURI)
	}
	for _, opt := range opts {
		opt.setValue(v)
	}
	return retrieveToken(ctx, shop, c, v)
}

func (c *Config) RefreshToken(ctx context.Context, shop, token string, opts ...AuthCodeOption) (*Token, error) {
	v := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {token},
	}
	for _, opt := range opts {
		opt.setValue(v)
	}
	return retrieveToken(ctx, shop, c, v)
}

func (c *Config) fixSite(shop string) string {
	return fmt.Sprintf("https://%s", shop)
}

func (c *Config) fixTokenUrl(shop string) string {
	return fmt.Sprintf("https://%s/%s", shop, strings.TrimPrefix(c.Endpoint.TokenURL, "/"))
}

func (c *Config) ValidShop(shop string) bool {
	domain := c.Domain
	if domain == "" {
		domain = DefaultDomain
	}
	shopRegexp := regexp.MustCompile("^[a-zA-Z0-9-]+." + domain + "$")
	return shopRegexp.MatchString(shop)
}

func retrieveToken(ctx context.Context, shop string, c *Config, v url.Values) (*Token, error) {
	tk, err := internal.RetrieveToken(ctx, c.ClientID, c.ClientSecret, c.fixTokenUrl(shop), v)
	if err != nil {
		if rErr, ok := err.(*internal.RetrieveError); ok {
			return nil, (*RetrieveError)(rErr)
		}
		return nil, err
	}
	return tokenFromInternal(tk), nil
}
