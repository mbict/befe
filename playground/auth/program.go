package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	. "github.com/mbict/befe/dsl"
	. "github.com/mbict/befe/dsl/http"
	"github.com/mbict/befe/dsl/templates"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	template    templates.HTMLTemplate
	whitelabels map[string]JSON
)

func Program() Expr {
	backend := ReverseProxy(FromEnvWithDefault("BACKEND_URI", "http://localhost:8082"))

	fs := FileServer("./dist/").
		AllowDirectoryListing(false).
		WhenNotFound(NotFound(), WriteResponse([]byte("the page does not exist 404 html here")))

	template = templates.New(
		templates.FromFile(`./index.gotmpl`),
	)

	whitelabels = loadWhitelabels("./whitelabels")

	d := Decision()
	{
		//unprotected endpoints, no need for JWT/JWK checking
		d.When(PathStartWith("/assets", "/img", "/favicon", "/manifest.json")).
			Then(fs)

		//account selector
		d.When(
			PathEquals("/"),
			HasCookie("ssid"),
		).Then(WriteResponse([]byte(`account selection app here`)))

		//redirect to login page if no session isset
		d.When(
			PathEquals("/"),
		).Then(TemporaryRedirect("/login")) /* make full url with domain */

		//redirect to check-session backend to verify session is valid
		//will handle redirects
		d.When(
			IsMethod(http.MethodGet, http.MethodOptions),
			PathEquals("/login"),
			HasCookie("ssid"),
		).Then(SetPath("/check-session"), backend)

		//backend auth calls
		d.When(
			IsMethod(http.MethodPost, http.MethodOptions),
			PathEquals("/oauth2/auth", "/oauth2/token", "/oauth2/revoke", "/oauth2/introspect", "/login", "/register", "/recover", "/verify"),
		).Then(backend)

		d.When(
			IsMethod(http.MethodGet, http.MethodOptions),
			PathEquals("/logout", "/userinfo", "/oauth2/auth", "/.well-known/openid-configuration", "/.well-known/jwks.json"),
		).Then(backend)

		d.When(PathEquals(loadPaths(`./paths/account.txt`)...)).
			Then(renderIndex("accounts"))

		//if no match was made we render the 404 page
		d.Else(pageNotFound())
	}

	cors := CORS().
		AllowedOrigins(FromEnvWithDefault("API_URI", "http://localhost/")).
		AllowAllMethods().
		AllowedHeaders("Authorization")

	return With(cors, d)

}

func renderIndex(entrypoint string) Expr {
	return With(
		Ok(),
		template.RenderTemplate("index.gotmpl",
			templates.WithData("whitelabelConfig", func(r *http.Request) interface{} {
				if config, ok := whitelabels[r.Host]; ok {
					return config
				}
				return whitelabels[`default`]
			}),
		),
	)
}

func pageNotFound() Expr {
	return With(
		NotFound(),
		WriteResponse([]byte(`page not found`)),
	)
}

func loadPaths(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("failed to load the path file `%s` with error :%v\n", filename, err.Error())
		return nil
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func loadWhitelabels(whitelabelPath string) map[string]JSON {
	result := make(map[string]JSON)
	err := filepath.Walk(whitelabelPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() /*|| filepath.Ext(path) != `.json`*/ {
			return nil
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var data JSON
		err = json.Unmarshal(contents, &data)
		if err != nil {
			return err
		}

		hostName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		result[hostName] = data

		return nil
	})

	if err != nil {
		fmt.Println("cannot unmarshal whitelabel file", err.Error())
		panic(err)
	}

	return result
}
