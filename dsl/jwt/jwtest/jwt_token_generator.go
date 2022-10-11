package jwtest

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"time"
)

type Claims map[string]interface{}

type JwtTokenGenerator interface {
	IsExpired() JwtTokenGenerator

	WithAudiences(...string) JwtTokenGenerator
	WithScopes(...string) JwtTokenGenerator
	WithSubject(string) JwtTokenGenerator
	WithIssuer(string) JwtTokenGenerator
	WithCustomClaim(string, interface{}) JwtTokenGenerator
	WithInvalidSigner() JwtTokenGenerator

	Generate() string
	GenerateBearer() string
}

func JwtGenerator() JwtTokenGenerator {
	return &jwtTokenGenerator{
		audiences:     nil,
		scopes:        nil,
		subject:       "",
		claims:        make(Claims),
		issuer:        "",
		isExpired:     false,
		invalidSigner: false,
	}
}

type jwtTokenGenerator struct {
	isExpired     bool
	invalidSigner bool
	audiences     []string
	scopes        []string
	subject       string
	issuer        string
	claims        Claims
}

func (j *jwtTokenGenerator) clone() *jwtTokenGenerator {
	claims := make(map[string]interface{}, len(j.claims))
	for k, v := range j.claims {
		claims[k] = v
	}

	jc := &jwtTokenGenerator{
		isExpired:     j.isExpired,
		invalidSigner: j.invalidSigner,
		audiences:     j.audiences,
		scopes:        j.scopes,
		subject:       j.subject,
		issuer:        j.issuer,
		claims:        claims,
	}

	return jc
}

func (j *jwtTokenGenerator) IsExpired() JwtTokenGenerator {
	j = j.clone()
	j.isExpired = true
	return j
}

func (j *jwtTokenGenerator) WithInvalidSigner() JwtTokenGenerator {
	j = j.clone()
	j.invalidSigner = true
	return j
}

func (j *jwtTokenGenerator) WithAudiences(audiences ...string) JwtTokenGenerator {
	j = j.clone()
	j.audiences = append(j.audiences, audiences...)
	return j
}

func (j *jwtTokenGenerator) WithScopes(scopes ...string) JwtTokenGenerator {
	j = j.clone()
	j.scopes = append(j.scopes, scopes...)
	return j
}

func (j *jwtTokenGenerator) WithSubject(subject string) JwtTokenGenerator {
	j = j.clone()
	j.subject = subject
	return j
}

func (j *jwtTokenGenerator) WithIssuer(issuer string) JwtTokenGenerator {
	j = j.clone()
	j.issuer = issuer
	return j
}

func (j *jwtTokenGenerator) WithCustomClaim(name string, value interface{}) JwtTokenGenerator {
	j = j.clone()
	j.claims[name] = value
	return j
}

func (j *jwtTokenGenerator) Generate() string {
	t := jwt.New()
	t.Set(jwt.SubjectKey, j.subject)

	for k, v := range j.claims {
		t.Set(k, v)
	}

	if j.isExpired == true {
		t.Set(jwt.IssuedAtKey, time.Now().Add(-180*time.Minute))
		t.Set(jwt.ExpirationKey, time.Now().Add(-60*time.Minute))
	} else {
		t.Set(jwt.IssuedAtKey, time.Now().Add(-60*time.Minute))
		t.Set(jwt.ExpirationKey, time.Now().Add(60*time.Minute))
	}

	if len(j.audiences) > 0 {
		t.Set(jwt.AudienceKey, j.audiences)
	}

	if len(j.scopes) > 0 {
		t.Set("scp", j.scopes)
	}

	if j.issuer != "" {
		t.Set(jwt.IssuerKey, j.issuer)
	}

	key := PrivateJwk()
	if j.invalidSigner == true {
		//we create a random rsa key
		rsaKey, _ := rsa.GenerateKey(rand.Reader, 2048)
		key, _ = jwk.New(rsaKey)
		key.Set(jwk.KeyIDKey, "invalid_signer")
		key.Set(jwk.AlgorithmKey, jwa.RS256)
	}

	signed, err := jwt.Sign(t, jwa.RS256, key)
	if err != nil {
		panic(fmt.Errorf("failed to sign token: %w", err))
	}

	return string(signed)
}

func (j *jwtTokenGenerator) GenerateBearer() string {
	return "Bearer " + j.Generate()
}
