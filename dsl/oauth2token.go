package dsl

import (
	"context"
	. "github.com/mbict/befe/expr"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"log"
	"net/http"
)

var oauth2TokenContextKey int

type OAuthClientToken interface {
	Action

	WhenDenied(actions ...Action) OAuthClientToken
	WhenError(actions ...Action) OAuthClientToken
	InjectToken() OAuthClientToken
}

func OAuthClientCredentials(clientId, clientSecret, tokenUrl string, scopes []string) OAuthClientToken {
	conf := clientcredentials.Config{
		ClientID:       clientId,
		ClientSecret:   clientSecret,
		Scopes:         scopes,
		TokenURL:       tokenUrl,
		EndpointParams: nil,
	}
	ctx := context.Background()

	return &oauthClientToken{
		tokenSource: conf.TokenSource(ctx),
	}
}

type oauthClientToken struct {
	tokenSource       oauth2.TokenSource
	injectAccessToken bool

	deniedActions Actions
	errorActions  Actions
}

func (o *oauthClientToken) InjectToken() OAuthClientToken {
	o.injectAccessToken = true
	return o
}

func (o *oauthClientToken) WhenDenied(actions ...Action) OAuthClientToken {
	o.deniedActions = actions
	return o
}

func (o *oauthClientToken) WhenError(actions ...Action) OAuthClientToken {
	o.errorActions = actions
	return o
}

func (o *oauthClientToken) BuildHandler(ctx context.Context, next Handler) Handler {
	deniedHandler := o.deniedActions.BuildHandler(ctx, nil)
	errorHandler := o.errorActions.BuildHandler(ctx, nil)
	return func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		t, err := o.tokenSource.Token()

		if err != nil {
			switch err.Error() {
			case "exp not satisfied":
				return deniedHandler(rw, r)
			default:
				log.Printf("unkown error decoding token :%v", err)
			}
			return errorHandler(rw, r)
		}

		//save the token in the context
		r = r.WithContext(context.WithValue(r.Context(), &oauth2TokenContextKey, t))

		//inject the auth token
		if o.injectAccessToken == true {
			t.SetAuthHeader(r)
		}

		if next == nil {
			return true, nil
		}
		return next(rw, r)
	}
}

func InjectOAuth2AccessToken() Action {
	return ActionFunc(func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		t, ok := r.Context().Value(&oauth2TokenContextKey).(*oauth2.Token)
		if ok {
			t.SetAuthHeader(r)
		}
		return true, nil
	})
}

func OAuth2AccessToken() Valuer {
	return func(r *http.Request) interface{} {
		t, ok := r.Context().Value(&oauth2TokenContextKey).(*oauth2.Token)
		if ok {
			return t.AccessToken
		}
		return ""
	}
}

func OAuth2AuthorizationAccessToken() Valuer {
	return func(r *http.Request) interface{} {
		t, ok := r.Context().Value(&oauth2TokenContextKey).(*oauth2.Token)
		if ok {
			return t.Type() + " " + t.AccessToken
		}
		return ""
	}
}
