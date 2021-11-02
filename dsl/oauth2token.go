package dsl

import (
	"context"
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
	deniedHandler := o.deniedActions.BuildHandler(ctx, emptyHandler)
	errorHandler := o.errorActions.BuildHandler(ctx, emptyHandler)
	return func(rw http.ResponseWriter, r *http.Request) {
		t, err := o.tokenSource.Token()

		if err != nil {
			switch err.Error() {
			case "exp not satisfied":
				deniedHandler(rw, r)
				return
			default:
				log.Printf("unkown error decoding jwt token :%v", err)
				//deniedHandler(rw, r)
			}
			errorHandler(rw, r)
			return
		}

		//save the token in the context
		r = r.WithContext(context.WithValue(r.Context(), &oauth2TokenContextKey, t))

		//inject the auth token
		if o.injectAccessToken == true {
			t.SetAuthHeader(r)
		}

		next(rw, r)
	}
}

func InjectOAuth2AccessToken() Action {
	return ActionBuilder(func(_ context.Context, next Handler) Handler {
		return func(rw http.ResponseWriter, r *http.Request) {
			t, ok := r.Context().Value(&oauth2TokenContextKey).(*oauth2.Token)
			if ok {
				t.SetAuthHeader(r)
			}

			if next != nil {
				next(rw, r)
			}
		}
	})
}
