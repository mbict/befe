package dsl

import (
	"context"
	"fmt"
	"github.com/mbict/befe/proxy"
	"github.com/mbict/befe/utils/bufferpool"
	"net/http"
	"net/url"
)

type reverseProxy struct {
	proxy *proxy.ReverseProxy
}

func (p *reverseProxy) BuildHandler(_ context.Context, next Handler) Handler {
	return func(rw http.ResponseWriter, r *http.Request) {
		transformer, _ := r.Context().Value(&transformContextKey).(TransformFunc)
		p.proxy.ServeHTTP(rw, r, proxy.ResponseModifier(transformer))

		if next != nil {
			next(rw, r)
		}
	}
}

func ReverseProxy(serviceUrl string) Action {
	u, err := url.Parse(serviceUrl)
	if err != nil {
		panic(fmt.Errorf("unable to parse url `%s` for reverse proxy: %w", serviceUrl, err))
	}

	p := proxy.NewSingleHostReverseProxy(u)
	p.BufferPool = bufferpool.New()

	return &reverseProxy{
		proxy: p,
	}
}
