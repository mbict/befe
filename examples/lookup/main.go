package main

import (
	. "github.com/mbict/befe/dsl"
)

// Program is an example where we use a lookup to check if the requested resource belongs to the authenticated user
func Program() Action {
	jwk := JwkToken(FromEnvWithDefault("JWK_URI", "http://localhost/.well-known/jwks.json")).
		WithExpiredCheck()

	objectBelongsToCustomer := HttpLookup("http://localhost/objects/{object_id}", WithParam("object_id", ValueFromQuery("object_id"))).
		Must(HaveResponseCode(200), JsonResponse(JsonHasValue("$.customer_id", ValueFromJwtClaim("sub")))).
		OnFail(NotFound(), WriteResponseBody([]byte(`object not found or failed`)))

	return With(jwk, objectBelongsToCustomer, WriteResponseBody([]byte(`object belongs to user`)))
}
