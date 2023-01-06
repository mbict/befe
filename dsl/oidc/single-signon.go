package oidc

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc/v3/oidc"
	. "github.com/mbict/befe/expr"
	"golang.org/x/oauth2"
	"net/http"
	"time"
)

type singleSignon struct {
	authority    string
	clientId     string
	clientSecret string
	redirectURL  string
}

func (s *singleSignon) BuildHandler(ctx context.Context, next Handler) Handler {

	provider, err := oidc.NewProvider(ctx, s.authority)
	if err != nil {
		panic(err)
	}

	// Configure an OpenID Connect aware OAuth2 client.
	oauth2Config := oauth2.Config{
		ClientID:     s.clientId,
		ClientSecret: s.clientSecret,
		RedirectURL:  s.redirectURL,

		// Discovery returns the OAuth2 endpoints.
		Endpoint: provider.Endpoint(),

		// "openid" is a required scope for OpenID Connect flows.
		Scopes: []string{ /*oidc.ScopeOpenID, "profile", "email"*/ },
	}

	var verifier = provider.Verifier(&oidc.Config{
		ClientID:                   s.clientId,
		SupportedSigningAlgs:       []string{"RS256"},
		SkipClientIDCheck:          true,
		SkipExpiryCheck:            false,
		SkipIssuerCheck:            true,
		Now:                        nil,
		InsecureSkipSignatureCheck: false,
	})

	unsetCookie := func() *http.Cookie {
		return &http.Cookie{
			Name:     "access_token",
			Value:    "",
			Path:     "/",
			Expires:  time.Unix(0, 0),
			HttpOnly: true,
		}

	}

	setCookie := func(accessToken string) *http.Cookie {
		return &http.Cookie{
			Name:  "access_token",
			Value: accessToken,
			Path:  "/",

			HttpOnly: true,
		}
	}

	return func(rw http.ResponseWriter, req *http.Request) (bool, error) {

		if cookie, err := req.Cookie("access_token"); err != http.ErrNoCookie {

			idToken, err := verifier.Verify(ctx, cookie.Value)
			if err != nil {
				http.SetCookie(rw, unsetCookie())

				// handle error
				switch e := err.(type) {
				case *oidc.TokenExpiredError:
					fmt.Println("EXPIRED TOKEN", e)

					fmt.Println("redirecting after expire", req.URL.Path)

					state := req.URL.Path + "12345678"

					http.Redirect(rw, req, oauth2Config.AuthCodeURL(state), http.StatusFound)

				default:
					fmt.Println("error verify in set cookie", e)
				}

				return false, nil
			}

			fmt.Println("checked token", string(idToken.AccessTokenHash))

			if req.URL.Path == `/auth-callback` {
				http.Redirect(rw, req, "/", http.StatusFound)
				return false, nil
			}

			return next(rw, req)
		}

		if req.URL.Path == `/auth-callback` {
			fmt.Println("validating")
			oauth2Token, err := oauth2Config.Exchange(ctx, req.URL.Query().Get("code"))
			if err != nil {
				// handle error
				fmt.Println("got error exchange", err)
				return false, nil
			}

			fmt.Println(oauth2Token)

			// Extract the ID Token from OAuth2 token.
			rawIDToken, ok := oauth2Token.Extra("access_token").(string)
			if !ok {
				// handle missing token
				fmt.Println("no id token")
				return false, nil
			}

			fmt.Println("access token", rawIDToken)

			// Parse and verify ID Token payload.
			idToken, err := verifier.Verify(ctx, rawIDToken)
			if err != nil {
				// handle error
				fmt.Println("error verify", err)
				return false, nil
			}

			fmt.Println("id token", string(idToken.AccessTokenHash))

			// Extract custom claims
			//var claims struct {
			//	Email    string `json:"email"`
			//	Verified bool   `json:"email_verified"`
			//}
			//if err := idToken.Claims(&claims); err != nil {
			//	// handle error
			//}

			http.SetCookie(rw, setCookie(rawIDToken))
			http.Redirect(rw, req, "/", http.StatusFound)

			return false, nil
			//return next(rw, req)
		} else {
			fmt.Println("redirecting", req.URL.Path)

			state := req.URL.Path + "12345678"

			http.Redirect(rw, req, oauth2Config.AuthCodeURL(state), http.StatusFound)
		}

		return false, nil
	}
}

func SingleSignon(authority string, clientId string, clientSecret string) Action {
	return &singleSignon{
		authority:    authority,
		clientId:     clientId,
		clientSecret: clientSecret,
		redirectURL:  "http://localhost:8083/auth-callback",
	}
}
