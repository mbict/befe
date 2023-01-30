package oidc

import (
	"errors"
	. "github.com/mbict/befe/expr"
	"net/http"
)

var (
	ErrNoToken         = errors.New("no token")        //no token in cookie
	ErrNoAccessToken   = errors.New("no access token") // no access token was returned in the authorize request
	ErrTokenExpired    = errors.New("token expired")
	ErrValidatingToken = errors.New("token validation failed")
)

var (
	jwtTokenContextKey int
	providerContextKey int
)

const authCallbackPath = `/auth/callback`

func WithCookieHttpOnly(httpOnly bool) ProviderOption {
	return func(provider *ssoProvider) {
		provider.cookieHttpOnly = httpOnly
	}
}

func WithCookiePath(path string) ProviderOption {
	return func(provider *ssoProvider) {
		provider.cookiePath = path
	}
}

type SSO interface {
	Action
	Condition

	WithSameIssuer() SSO
	WithAudience(audiences ...string) SSO
	WithClaim(name string, values ...string) SSO

	WhenExpired(...Action) SSO
	WhenDenied(...Action) SSO
	WhenInvalidToken(...Action) SSO
	WhenNoToken(...Action) SSO
}

// AuthTokenRedirect will redirect the user to the openid endpoint to get authenticated and obtain an access token
func AuthTokenRedirect() Action {
	return ActionFunc(func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		//if we already tried to fetch the provider, we can do it here
		provider, ok := r.Context().Value(&providerContextKey).(*ssoProvider)
		if !ok {
			//we could not retrieve the oidc provider config
			return false, errors.New("could not retrieve oidc provider config")
		}

		provider.redirect(rw, r)
		return false, nil
	})
}
