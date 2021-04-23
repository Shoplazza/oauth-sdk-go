package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/context/ctxhttp"
)

type StoreInfo struct {
	ID     int
	Locale string
	Name   string
}

type Token struct {
	// AccessToken is the token that authorizes and authenticates the requests.
	AccessToken string

	// TokenType is the type of token. The Type method returns either this or "Bearer", the default.
	TokenType string

	RefreshToken string

	ExpiresAt time.Time
	Scopes    []string
	StoreInfo StoreInfo
}

// tokenJSON is the struct representing the HTTP response from OAuth2 providers returning a token in JSON form.
type tokenJSON struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	Scope        string `json:"scope"`
	CreatedAt    string `json:"created_at"`
	ExpiresAt    string `json:"expires_at"`
	StoreID      string `json:"store_id"`
	StoreName    string `json:"store_name"`
	Locale       string `json:"locale"`
}

type expirationTime int32

func (e *expirationTime) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		return nil
	}
	var n json.Number
	err := json.Unmarshal(b, &n)
	if err != nil {
		return err
	}
	i, err := n.Int64()
	if err != nil {
		return err
	}
	if i > math.MaxInt32 {
		i = math.MaxInt32
	}
	*e = expirationTime(i)
	return nil
}

// RegisterBrokenAuthHeaderProvider previously did something. It is now a no-op.
//
// Deprecated: this function no longer does anything. Caller code that
// wants to avoid potential extra HTTP requests made during
// auto-probing of the provider's auth style should set
// Endpoint.AuthStyle.
func RegisterBrokenAuthHeaderProvider(tokenURL string) {}

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
	if err != nil && needsAuthStyleProbe {
		authStyle = AuthStyleInParams // the second way we'll try
		req, _ = newTokenRequest(tokenURL, clientID, clientSecret, v, authStyle)
		token, err = doTokenRoundTrip(ctx, req)
	}
	if needsAuthStyleProbe && err == nil {
		setAuthStyle(tokenURL, authStyle)
	}
	// Don't overwrite `RefreshToken` with an empty value
	// if this was a token refreshing request.
	if token != nil && token.RefreshToken == "" {
		token.RefreshToken = v.Get("refresh_token")
	}
	return token, err
}

func doTokenRoundTrip(ctx context.Context, req *http.Request) (*Token, error) {
	r, err := ctxhttp.Do(ctx, ContextClient(ctx), req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1<<20))
	r.Body.Close()
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
			Raw:          vals,
		}
		e := vals.Get("expires_in")
		expires, _ := strconv.Atoi(e)
		if expires != 0 {
			token.Expiry = time.Now().Add(time.Duration(expires) * time.Second)
		}
	default:
		var tj tokenJSON
		if err = json.Unmarshal(body, &tj); err != nil {
			return nil, err
		}
		token = &Token{
			AccessToken:  tj.AccessToken,
			TokenType:    tj.TokenType,
			RefreshToken: tj.RefreshToken,
			Expiry:       tj.expiry(),
			Raw:          make(map[string]interface{}),
		}
		json.Unmarshal(body, &token.Raw) // no error checks for optional fields
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
	return fmt.Sprintf("oauth2: cannot fetch token: %v\nResponse: %s", r.Response.Status, r.Body)
}
