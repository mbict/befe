package oidc

import (
	"context"
	. "github.com/mbict/befe/expr"
	"net/http"
	"sync"
)

func SingleSignOn(authority string, clientId string, clientSecret string, redirectUrl string, options ...ProviderOption) SSO {
	return &singleSignon{
		authority:       authority,
		clientId:        clientId,
		clientSecret:    clientSecret,
		redirectURL:     redirectUrl + authCallbackPath,
		providerOptions: options,
	}
}

type singleSignon struct {
	provider        *ssoProvider
	providerOptions []ProviderOption
	mu              sync.Mutex

	authority    string
	clientId     string
	clientSecret string
	redirectURL  string

	expiredActions      Actions
	deniedActions       Actions
	invalidTokenActions Actions
	noTokenActions      Actions
}

func (s *singleSignon) WithSameIssuer() SSO {
	//TODO implement me
	return s
}

func (s *singleSignon) WithAudience(audiences ...string) SSO {
	//TODO implement me
	panic("implement me sso with audience")
}

func (s *singleSignon) WithExpiredCheck() SSO {
	//TODO implement me
	panic("implement me sso with with expired check")
}

func (s *singleSignon) WithClaim(name string, values ...string) SSO {
	//TODO implement me
	panic("implement me sso with claiom")
}

func (s *singleSignon) WhenExpired(action ...Action) SSO {
	s.expiredActions = append(s.expiredActions, action...)
	return s
}

func (s *singleSignon) WhenDenied(action ...Action) SSO {
	s.invalidTokenActions = append(s.invalidTokenActions, action...)
	return s
}

func (s *singleSignon) WhenInvalidToken(action ...Action) SSO {
	s.invalidTokenActions = append(s.invalidTokenActions, action...)
	return s
}

func (s *singleSignon) WhenNoToken(action ...Action) SSO {
	s.noTokenActions = append(s.noTokenActions, action...)
	return s
}

func (s *singleSignon) getProvider(ctx context.Context) *ssoProvider {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.provider != nil {
		return s.provider
	}

	provider, err := newProvider(ctx, s.authority, s.clientId, s.clientSecret, s.redirectURL, s.providerOptions...)
	if err != nil {
		panic(err)
	}

	s.provider = provider
	return provider
}

func (s *singleSignon) BuildCondition(ctx context.Context) ConditionFunc {
	provider := s.getProvider(ctx)
	return func(req *http.Request) bool {
		//check cookie
		idToken, err := provider.checkCookie(req)
		if err != nil {
			return false
		}

		//we do nothing here yet, we can check more here later
		_ = idToken

		return true
	}
}

func (s *singleSignon) BuildHandler(ctx context.Context, next Handler) Handler {
	provider := s.getProvider(ctx)

	var onNoToken Handler
	if len(s.noTokenActions) > 0 {
		onNoToken = s.noTokenActions.BuildHandler(ctx, nil)
	}

	var onInvalidToken Handler
	if len(s.invalidTokenActions) > 0 {
		onInvalidToken = s.invalidTokenActions.BuildHandler(ctx, nil)
	}

	var onExpiredToken Handler
	if len(s.expiredActions) > 0 {
		onExpiredToken = s.expiredActions.BuildHandler(ctx, nil)
	}

	return func(rw http.ResponseWriter, req *http.Request) (bool, error) {
		req = req.WithContext(context.WithValue(req.Context(), &providerContextKey, provider))
		return provider.handleRequest(rw, req, next, onExpiredToken, onNoToken, onInvalidToken)
	}
}
