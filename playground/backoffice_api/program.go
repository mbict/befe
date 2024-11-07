package main

import (
	. "github.com/mbict/befe/dsl"
	. "github.com/mbict/befe/dsl/http"
	"github.com/mbict/befe/dsl/jwt"
)

var (
	//checks based on jwt claims
	hasAllowedRole = Or(jwt.HasJwtClaim("urp", "picker"), jwt.HasJwtClaim("urp", "owner"), jwt.HasJwtClaim("urp", "administrator"))
	//isCustomer      = jwt.HasJwtClaim("urp", "customer")

	coreApiClient HttpClient
)

func Program() Expr {

	//core api
	coreApiBasePath := PrependPath("/stalling/{stallingId}", jwt.ParamFromClaim("stallingId", "aud"))
	//coreApiBasePath := PrependPath("/stalling/{stallingId}", ParamString("stallingId", "5abbb6e6-5638-4e3a-814b-3a14aaaa8ddf"))
	coreApi := With(
		coreApiBasePath,
		ReverseProxy(
			FromEnvWithDefault("CORE_API_URI", "https://api.dev.stalling.app/core"),
			//FromEnvWithDefault("CORE_API_URI", "http://localhost:8084/api"),
			WithHostFromServiceUrl(), //only here for development purposes as i'm using a different url than the upstream can handle
		),
	)

	//client used for enriching
	coreApiClient = Client(FromEnvWithDefault("CORE_API_URI", "https://api.dev.stalling.app/core")) //reverse proxy

	//core api
	authApi := ReverseProxy(
		FromEnvWithDefault("AUTH_API_URI", "https://api.dev.stalling.app/auth"),
		WithHostFromServiceUrl(), //only here for development purposes as i'm using a different url than the upstream can handle
	)
	/* todo: tighten this! */
	cors := CORS().
		AllowAll().
		ExposedHeaders("Location", "Content-Type").
		AllowedHeaders("Authorization", "Content-Type").
		MaxAge(3600)

	jwk := jwt.JwkToken(FromEnvWithDefault("JWK_URI", "https://accounts.dev.stalling.app/.well-known/jwks.json")).
		WithExpiredCheck().
		/* todo: dynamicly load all the audiences from the whitelabel configs */
		//WithAudience("api://stalling/*").
		WhenExpired(JsonContentType(), Deny(), WriteJson(JSON{"error": "expired_token"})).
		WhenDenied(JsonContentType(), Deny(), WriteJson(JSON{"error": "invalid_token"}))

	r := Router().With(cors, jwk).
		OnNotFound(NotFound(), WriteJson(JSON{"error": "not_found"}))

	//stalling account

	//locations
	r.Get("/locations").Then(coreApi)
	r.Get("/locations/{locationId}").Then(coreApi)
	r.Post("/locations").Then(coreApi)
	r.Patch("/locations/{locationId}").Then(coreApi)
	r.Put("/locations/{locationId}:rename").Then(coreApi)
	r.Put("/locations/{locationId}:markAsPickupPoint").Then(coreApi)
	r.Put("/locations/{locationId}:markAsParkingSpace").Then(coreApi)
	r.Put("/locations/{locationId}:changeImage").Then(coreApi)
	r.Put("/locations/{locationId}:removeImage").Then(coreApi)
	r.Delete("/locations/{locationId}").Then(coreApi)

	//parkings
	r.Get("/parkings").Then(coreApi, locationsEnricher())
	r.Get("/parkings/{parkingId}").Then(coreApi, locationEnricher())
	r.Post("/parkings").Then(coreApi)
	r.Delete("/parkings/{parkingId}").Then(coreApi)

	//objects
	r.Get("/objects").Then(coreApi)
	r.Get("/objects/{objectId}").Then(coreApi)
	r.Post("/objects").Then(coreApi)
	r.Delete("/objects/{objectId}").Then(coreApi)

	//customers
	r.Get("/customers").Then(coreApi)
	r.Get("/customers/{customerId}").Then(coreApi)
	r.Post("/customers").Then(coreApi)
	r.Delete("/customers/{customerId}").Then(coreApi)

	//identifiers
	r.Get("/identifiers").Then(coreApi)
	r.Post("/identifiers").Then(coreApi)
	r.Get("/identifiers:findAvailableIdentifiers").Then(coreApi)
	r.Delete("/identifiers/{identifierId}").Then(coreApi)

	//users
	authUsersBasePath := SetPath("/tenants/{stallingId}/users", jwt.ParamFromClaim("stallingId", "aud"))
	r.Get("/users").Then(authUsersBasePath, authApi)

	return r
	//
	//return With(
	//	cors,
	//	//jwk,
	//	Decision().
	//		When(hasAllowedRole).Then(backofficeApi()).
	//		Else(Deny(), WriteJson(JSON{"error": "invalid_token"})),
	//)

}

func backofficeApi() Expr {

	return WriteJson(JSON{"ok": "data"})
}

func locationsEnricher() Expr {
	return coreApiClient.Get( //this api is paginated by offset, so we need to visit all pages to get a full result
		UrlBuilder("/stalling/{stallingId}/locations?location_ids={locationIds}&size={size}",
			ParamValue("size", 1000),                                 //we do amaximum fixed size
			jwt.ParamFromClaim("stallingId", "aud"),                  //stallind is fetched from jwt
			ParamFromJsonPath("locationIds", "$.data.*.location_id"), //gather all the location ids from the result
		)).
		OnSuccess(
			ResultMerger().
				Target("data.*.location").
				Source(ValueFromJsonPath("$.data.*")).
				Matcher("$.id", "$.location_id"),
		).
		OnFailure(
			InternalServerError(),
			Stop(), //we will not continue we will hard stop here
		)
}

func locationEnricher() Expr {

	HasQuery("include")

	return coreApiClient.Get( //this api is paginated by offset, so we need to visit all pages to get a full result
		UrlBuilder("/stalling/{stallingId}/locations/{locationId}",
			jwt.ParamFromClaim("stallingId", "aud"),          //stallingId is fetched from jwt
			ParamFromJsonPath("locationId", "$.location_id"), //get the location id
		)).
		OnSuccess(
			ResultMerger().
				Target("location").
				Source(ValueFromJsonPath("$.*")).
				Matcher("$.id", "$.location_id"),
		).
		OnFailure(
			InternalServerError(),
			Stop(), //we will not continue we will hard stop here
		)
}
