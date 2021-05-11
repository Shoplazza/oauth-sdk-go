package oauth2

import (
	"fmt"
	"net/http"
	"net/url"
	"oauth2/internal"
	"strconv"
	"strings"
	"time"
)

const expiryDelta = 10 * time.Second

type StoreInfo struct {
	ID     int
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
	raw          interface{}
}

// Type returns t.TokenType if non-empty, else "Bearer".
func (t *Token) Type() string {
	if strings.EqualFold(t.TokenType, "bearer") {
		return "Bearer"
	}
	if strings.EqualFold(t.TokenType, "mac") {
		return "MAC"
	}
	if strings.EqualFold(t.TokenType, "basic") {
		return "Basic"
	}
	if t.TokenType != "" {
		return t.TokenType
	}
	return "Bearer"
}

// tokenFromInternal maps an *internal.Token struct into
// a *Token struct.
func tokenFromInternal(t *internal.Token) *Token {
	if t == nil {
		return nil
	}
	token := &Token{
		AccessToken:  t.AccessToken,
		TokenType:    t.TokenType,
		RefreshToken: t.RefreshToken,
		ExpiresAt:    t.ExpiresAt,
		raw:          t.Raw,
	}
	if t.StoreInfo != nil {
		token.StoreInfo = &StoreInfo{
			ID:     t.StoreInfo.ID,
			Name:   t.StoreInfo.Name,
			Locale: t.StoreInfo.Locale,
		}
	}
	return token
}

// Extra returns an extra field.
// Extra fields are key-value pairs returned by the server as a
// part of the token retrieval response.
func (t *Token) Extra(key string) interface{} {
	if raw, ok := t.raw.(map[string]interface{}); ok {
		return raw[key]
	}

	vals, ok := t.raw.(url.Values)
	if !ok {
		return nil
	}

	v := vals.Get(key)
	switch s := strings.TrimSpace(v); strings.Count(s, ".") {
	case 0: // Contains no "."; try to parse as int
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			return i
		}
	case 1: // Contains a single "."; try to parse as float
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return f
		}
	}

	return v
}

// timeNow is time.Now but pulled out as a variable for tests.
var timeNow = time.Now

// expired reports whether the token is expired.
// t must be non-nil.
func (t *Token) expired() bool {
	if t.ExpiresAt.IsZero() {
		return false
	}
	return t.ExpiresAt.Round(0).Add(-expiryDelta).Before(timeNow())
}

// Valid reports whether t is non-nil, has an AccessToken, and is not expired.
func (t *Token) Valid() bool {
	return t != nil && t.AccessToken != "" && !t.expired()
}

// RetrieveError is the error returned when the token endpoint returns a
// non-2XX HTTP status code.
type RetrieveError struct {
	Response *http.Response
	// Body is the body that was consumed by reading Response.Body.
	// It may be truncated.
	Body []byte
}

func (r *RetrieveError) Error() string {
	return fmt.Sprintf("oauth2: cannot fetch token: %v\nResponse: %s", r.Response.Status, r.Body)
}
