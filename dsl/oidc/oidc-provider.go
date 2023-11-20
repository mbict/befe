package oidc

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/mbict/befe/dsl/jwt/jwtoken"
	"github.com/mbict/befe/expr"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

type ProviderOption func(provider *ssoProvider)

func newProvider(ctx context.Context, authority, clientId, clientSecret, redirectUrl string, options ...ProviderOption) (*ssoProvider, error) {
	provider, err := oidc.NewProvider(ctx, authority)
	if err != nil {
		return nil, err
	}

	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config := oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUrl,
		Endpoint:     provider.Endpoint(),                                  // Discovery returns the OAuth2 endpoints.
		Scopes:       []string{ /*oidc.ScopeOpenID, "profile", "email"*/ }, // "openid" is a required scope for OpenID Connect flows.
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID:                   oauth2Config.ClientID,
		SupportedSigningAlgs:       []string{"RS256"},
		SkipClientIDCheck:          true,
		SkipExpiryCheck:            false,
		SkipIssuerCheck:            true,
		Now:                        nil,
		InsecureSkipSignatureCheck: false,
	})

	p := &ssoProvider{
		provider:       provider,
		oauth2Config:   oauth2Config,
		verifier:       verifier,
		cookiePath:     "/",
		cookieHttpOnly: true,
	}

	for _, option := range options {
		option(p)
	}

	return p, nil
}

type ssoProvider struct {
	oauth2Config oauth2.Config
	provider     *oidc.Provider
	verifier     *oidc.IDTokenVerifier
	redirectUrl  string

	cookiePath     string
	cookieHttpOnly bool
}

func (p *ssoProvider) redirect(rw http.ResponseWriter, req *http.Request) {
	state := encodeUrlPathState(req.URL.Path)
	http.Redirect(rw, req, p.oauth2Config.AuthCodeURL(state), http.StatusFound)
}

func (p *ssoProvider) checkCookie(req *http.Request) (*oidc.IDToken, error) {
	cookie, err := req.Cookie("access_token")
	if err != http.ErrNoCookie && cookie.Value != "" {
		return p.verifyToken(req.Context(), cookie.Value)
	}
	return nil, ErrNoToken
}

func (p *ssoProvider) verifyToken(ctx context.Context, token string) (*oidc.IDToken, error) {
	idToken, err := p.verifier.Verify(ctx, token)
	if err != nil {
		// handle error
		switch e := err.(type) {
		case *oidc.TokenExpiredError:
			return nil, ErrTokenExpired

		default:
			//log these errors
			fmt.Println("error verify token:", e)
		}
		return nil, ErrValidatingToken
	}
	return idToken, nil
}

func (p *ssoProvider) handleRequest(rw http.ResponseWriter, req *http.Request, next, onExpiredTokenHandler, onNoTokenHandler, onInvalidTokenHandler expr.Handler) (bool, error) {
	//handle auth callback
	if req.URL.Path == authCallbackPath {
		jwtToken, idToken, err := p.handleAuthCallback(req)
		if err != nil {
			if onInvalidTokenHandler != nil {
				return onInvalidTokenHandler(rw, req)
			}
			return false, err
		}

		req = req.WithContext(jwtoken.ToContext(req.Context(), jwtoken.New(idToken)))
		// Extract custom claims
		//var claims struct {
		//	Email    string `json:"email"`
		//	Verified bool   `json:"email_verified"`
		//}
		//if err := idToken.Claims(&claims); err != nil {
		//	// handle error
		//}

		//set the cookie
		http.SetCookie(rw, &http.Cookie{
			Name:     "access_token",
			Value:    jwtToken,
			Path:     p.cookiePath,
			HttpOnly: p.cookieHttpOnly,
		})

		//redirect to the original page
		http.Redirect(rw, req, decodeUrlPathState(req.URL.Query().Get("state")), http.StatusFound)

		return false, nil
	}

	//we have cookie, lets validate and proceed if applicable
	idToken, err := p.checkCookie(req)
	if err != nil {
		http.SetCookie(rw, &http.Cookie{
			Name:     "access_token",
			Value:    "",
			Path:     p.cookiePath,
			Expires:  time.Unix(0, 0),
			HttpOnly: p.cookieHttpOnly,
		})
		switch err {
		case ErrNoToken:
			if onNoTokenHandler == nil {
				//redirect to oauth
				p.redirect(rw, req)
				return false, nil
			}
			_, err = onNoTokenHandler(rw, req)
			return false, err

		case ErrTokenExpired:
			if onExpiredTokenHandler == nil {
				//redirect to oauth
				p.redirect(rw, req)
				return false, nil
			}
			_, err := onExpiredTokenHandler(rw, req)
			return false, err

		case ErrValidatingToken:
			if onInvalidTokenHandler != nil {
				_, _ = onInvalidTokenHandler(rw, req)
			}
			return false, err

		default:
			fmt.Println("unhandled error worth to log it", err)
			if onInvalidTokenHandler != nil {
				_, _ = onInvalidTokenHandler(rw, req)
			}
			return false, err
		}
	}
	req = req.WithContext(jwtoken.ToContext(req.Context(), jwtoken.New(idToken)))

	//we should navigate away from the callback auth page, to avoid
	if req.URL.Path == authCallbackPath {
		http.Redirect(rw, req, "/", http.StatusFound)
		return false, nil
	}

	//execute the normal flow
	return next(rw, req)
}

func (p *ssoProvider) handleAuthCallback(req *http.Request) (string, *oidc.IDToken, error) {
	oauth2Token, err := p.oauth2Config.Exchange(req.Context(), req.URL.Query().Get("code"))
	if err != nil { //something went wrong during code/token exchange
		fmt.Println("got error exchange", err)
		return "", nil, err //todo return valid error
	}

	// Extract the ID Token from OAuth2 token.
	rawIDToken, ok := oauth2Token.Extra("access_token").(string)
	if !ok { //no access token to store
		return "", nil, ErrNoAccessToken
	}

	idToken, err := p.verifyToken(req.Context(), rawIDToken)
	if err != nil {
		return "", nil, err
	}
	return rawIDToken, idToken, nil
}
