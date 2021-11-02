package dsltest

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	"time"
)

var _privateKey *rsa.PrivateKey

func init() {
	//setup ceritificates to use with the test tokens
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(fmt.Errorf("failed to generate new RSA private key: %w", err))
		return
	}
	_privateKey = key

}

//Creates a valid jwk keyset to validate jwk tokens against
func JwkKeySet() []byte {
	set := jwk.NewSet()
	set.Add(jwkPublicKey())

	res, _ := json.MarshalIndent(set, "", "  ")
	return res
}

func JwkKeyS() jwk.Set {
	set := jwk.NewSet()
	set.Add(jwkPublicKey())

	return set
}

func jwkPrivateKey() jwk.Key {
	key, err := jwk.New(_privateKey)
	if err != nil {
		panic(fmt.Errorf("failed to create symmetric key: %w", err))
	}
	key.Set(jwk.KeyIDKey, "test")
	key.Set(jwk.AlgorithmKey, jwa.RS256)
	return key
}

func jwkPublicKey() jwk.Key {
	key, err := jwk.New(_privateKey.PublicKey)
	if err != nil {
		panic(fmt.Errorf("failed to create symmetric key: %w", err))
	}
	key.Set(jwk.KeyIDKey, "test")
	key.Set(jwk.AlgorithmKey, jwa.RS256)
	return key
}

type JwtTokenGenerator interface {
	IsExpired() JwtTokenGenerator

	WithAudiences(...string) JwtTokenGenerator
	WithScopes(...string) JwtTokenGenerator
	WithSubject(string) JwtTokenGenerator
	WithCustomClaim(string, interface{}) JwtTokenGenerator
	WithInvalidSigner() JwtTokenGenerator

	Generate() string
}

func JwtGenerator() JwtTokenGenerator {
	return &jwtTokenGenerator{
		audiences:     nil,
		scopes:        nil,
		subject:       "",
		claims:        make(map[string]interface{}),
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
	claims        map[string]interface{}
}

func (j *jwtTokenGenerator) IsExpired() JwtTokenGenerator {
	j.isExpired = true
	return j
}

func (j *jwtTokenGenerator) WithInvalidSigner() JwtTokenGenerator {
	j.invalidSigner = true
	return j
}

func (j *jwtTokenGenerator) WithAudiences(audiences ...string) JwtTokenGenerator {
	j.audiences = append(j.audiences, audiences...)
	return j
}

func (j *jwtTokenGenerator) WithScopes(scopes ...string) JwtTokenGenerator {
	j.scopes = append(j.scopes, scopes...)
	return j
}

func (j *jwtTokenGenerator) WithSubject(subject string) JwtTokenGenerator {
	j.subject = subject
	return j
}

func (j *jwtTokenGenerator) WithCustomClaim(name string, value interface{}) JwtTokenGenerator {
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

	key := jwkPrivateKey()
	if j.invalidSigner == true {
		//we create a random rsa key
		rsaKey, _ := rsa.GenerateKey(rand.Reader, 2048)
		key, _ = jwk.New(rsaKey)
		key.Set(jwk.KeyIDKey, "test")
		key.Set(jwk.AlgorithmKey, jwa.RS256)
	}

	signed, err := jwt.Sign(t, jwa.RS256, key)
	if err != nil {
		panic(fmt.Errorf("failed to sign token: %w", err))
	}

	return string(signed)
}

func ValidJwtToken(sub string) string {
	return JwtGenerator().
		WithSubject(sub).
		Generate()
}

func InvalidSignedJwtToken(sub string) string {
	return JwtGenerator().
		WithInvalidSigner().
		WithSubject(sub).
		Generate()
}

func ExpiredJwtToken(sub string) string {
	return JwtGenerator().
		IsExpired().
		WithSubject(sub).
		Generate()
}
