package oidc

import (
	"context"
	. "github.com/mbict/befe/expr"
	"net/http"
	"sync"
)

// DynamicSingleSignOn uses a dynamic authority and conifguration that can be created on the fly with Valuers
// authority is created and checked every request if there is a know configuration
// clientId, clientSecret and redirectUrl is only used when creating a new SSO config if the authority is not known yet.
func DynamicSingleSignOn(authority Valuer, clientId Valuer, clientSecret Valuer, redirectUrl Valuer) SSO {
	return &dynamicSSO{
		providers:          make(map[string]*ssoProvider),
		authorityValuer:    authority,
		clientIdValuer:     clientId,
		clientSecretValuer: clientSecret,
		redirectUrlValuer:  redirectUrl,
	}
}

type dynamicSSO struct {
	providers map[string]*ssoProvider
	mu        sync.RWMutex

	authorityValuer    Valuer
	clientIdValuer     Valuer
	clientSecretValuer Valuer
	redirectUrlValuer  Valuer

	expiredActions      Actions
	deniedActions       Actions
	invalidTokenActions Actions
	noTokenActions      Actions
}

func (s *dynamicSSO) WithSameIssuer() SSO {
	//TODO implement me
	return s
}

func (s *dynamicSSO) WithAudience(audiences ...string) SSO {
	//TODO implement me
	return s
}

func (s *dynamicSSO) WithClaim(name string, values ...string) SSO {
	//TODO implement me
	return s
}

func (s *dynamicSSO) WhenExpired(action ...Action) SSO {
	s.expiredActions = append(s.expiredActions, action...)
	return s
}

func (s *dynamicSSO) WhenDenied(action ...Action) SSO {
	s.deniedActions = append(s.deniedActions, action...)
	return s
}

func (s *dynamicSSO) WhenInvalidToken(action ...Action) SSO {
	s.invalidTokenActions = append(s.invalidTokenActions, action...)
	return s
}

func (s *dynamicSSO) WhenNoToken(action ...Action) SSO {
	s.noTokenActions = append(s.noTokenActions, action...)
	return s
}

func (s *dynamicSSO) BuildCondition(_ context.Context) ConditionFunc {
	return func(req *http.Request) bool {
		//get provider context
		provider, err := s.getProvider(req)
		if err != nil {
			return false
		}

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

func (s *dynamicSSO) BuildHandler(ctx context.Context, next Handler) Handler {
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
		provider, err := s.getProvider(req)
		if err != nil {
			return false, err
		}
		req = req.WithContext(context.WithValue(req.Context(), &providerContextKey, provider))

		return provider.handleRequest(rw, req, next, onExpiredToken, onNoToken, onInvalidToken)
	}
}

func (s *dynamicSSO) getProvider(r *http.Request) (*ssoProvider, error) {
	authority := s.authorityValuer(r).(string)

	s.mu.RLock()
	if provider, ok := s.providers[authority]; ok {
		s.mu.RUnlock()
		return provider, nil
	}
	s.mu.RUnlock()

	//if not found we try to create the provider
	s.mu.Lock()
	defer s.mu.Unlock()

	provider, err := newProvider(r.Context(),
		authority,
		s.clientIdValuer(r).(string),
		s.clientSecretValuer(r).(string),
		"http://"+s.redirectUrlValuer(r).(string)+authCallbackPath,
	)

	if err != nil {
		panic(err)
		return nil, err
	}

	s.providers[authority] = provider
	return provider, nil
}
