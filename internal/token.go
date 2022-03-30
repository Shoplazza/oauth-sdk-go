package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

type StoreInfo struct {
	ID     string
	Locale string
	Name   string
}

type Token struct {
	// AccessToken is the token that authorizes and authenticates the requests.
	AccessToken string
	// TokenType is the type of token. The Type method returns either this or "Bearer", the default.
	TokenType    string
	RefreshToken string
	ExpiresAt    time.Time
	StoreInfo    *StoreInfo
	Raw          interface{}
}

// tokenJSON is the struct representing the HTTP response from OAuth2 providers returning a token in JSON form.
type tokenJSON struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
	CreatedAt    int64  `json:"created_at"`
	ExpiresAt    int64  `json:"expires_at"`
	StoreID      string `json:"store_id"`
	StoreName    string `json:"store_name"`
	Locale       string `json:"locale"`
}

func newTokenRequest(tokenURL, clientID, clientSecret string, v url.Values) (*http.Request, error) {
	v = cloneURLValues(v)
	if clientID != "" {
		v.Set("client_id", clientID)
	}
	if clientSecret != "" {
		v.Set("client_secret", clientSecret)
	}
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("user-agent", "Oauth-SDK-Go/v1.0.4")
	return req, nil
}

func cloneURLValues(v url.Values) url.Values {
	v2 := make(url.Values, len(v))
	for k, vv := range v {
		v2[k] = append([]string(nil), vv...)
	}
	return v2
}

func RetrieveToken(ctx context.Context, clientID, clientSecret, tokenURL string, v url.Values) (*Token, error) {
	req, err := newTokenRequest(tokenURL, clientID, clientSecret, v)
	if err != nil {
		return nil, err
	}
	token, err := doTokenRoundTrip(ctx, req)
	return token, err
}

func doTokenRoundTrip(ctx context.Context, req *http.Request) (*Token, error) {
	r, err := ctxhttp.Do(ctx, http.DefaultClient, req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<20))
	_ = r.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("oauth2: cannot fetch token: %v", err)
	}
	if code := r.StatusCode; code < 200 || code > 299 {
		return nil, &RetrieveError{
			Response: r,
			Body:     body,
		}
	}

	var token *Token
	content, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	switch content {
	case "application/x-www-form-urlencoded", "text/plain":
		vals, err := url.ParseQuery(string(body))
		if err != nil {
			return nil, err
		}

		token = &Token{
			AccessToken:  vals.Get("access_token"),
			TokenType:    vals.Get("token_type"),
			RefreshToken: vals.Get("refresh_token"),
			StoreInfo: &StoreInfo{
				ID:     vals.Get("store_id"),
				Locale: vals.Get("locale"),
				Name:   vals.Get("store_name"),
			},
			Raw: vals,
		}
		expiresAt, _ := strconv.ParseInt(vals.Get("expires_at"), 10, 64)
		token.ExpiresAt = time.Unix(expiresAt, 0)
	default:
		var tj tokenJSON
		if err = json.Unmarshal(body, &tj); err != nil {
			return nil, err
		}
		token = &Token{
			AccessToken:  tj.AccessToken,
			TokenType:    tj.TokenType,
			RefreshToken: tj.RefreshToken,
			ExpiresAt:    time.Unix(tj.ExpiresAt, 0),
			Raw:          make(map[string]interface{}),
			StoreInfo: &StoreInfo{
				ID:     tj.StoreID,
				Name:   tj.StoreName,
				Locale: tj.Locale,
			},
		}
		_ = json.Unmarshal(body, &token.Raw) // no error checks for optional fields
	}
	if token.AccessToken == "" {
		return nil, errors.New("oauth2: server response missing access_token")
	}
	return token, nil
}

type RetrieveError struct {
	Response *http.Response
	Body     []byte
}

func (r *RetrieveError) Error() string {
	return fmt.Sprintf("Cannot fetch token: %v\nResponse: %s", r.Response.Status, r.Body)
}
