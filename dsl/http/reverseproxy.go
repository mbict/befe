package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	. "github.com/mbict/befe/expr"
	"github.com/mbict/befe/utils/bufferpool"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

type reverseProxy struct {
	proxy                 *httputil.ReverseProxy
	targetUrl             *url.URL
	useHostFromServiceUrl bool
}

func (p *reverseProxy) BuildHandler(_ context.Context, next Handler) Handler {
	rproxy := httputil.NewSingleHostReverseProxy(p.targetUrl)
	rproxy.BufferPool = bufferpool.New()

	//normally the host is forwarded, but in some situations like local development
	//you want to overwrite the host and use the host from the targeturl
	if p.useHostFromServiceUrl {
		currentDirector := rproxy.Director
		rproxy.Director = func(request *http.Request) {
			currentDirector(request)
			request.Host = request.URL.Host
		}
	}

	if next != nil {
		rproxy.ModifyResponse = func(response *http.Response) error {
			//try to decode body to a processable type
			data, err := DecodeResponse(response)
			if err != nil {
				return err
			}

			req := StoreResultBucket(response.Request, "", data, response, nil)

			rw := &proxyResponseWriter{}

			if _, err := next(rw, req); err != nil {
				return err
			}

			bucket := GetResultBucket(req.Context())

			fmt.Println("---->", bucket.Data, data)

			/* @todo hardcoded json encode to, dynamic encode

			rw.Header().Get(`Content-Type`)

			if rw.Len() > 0 {

			} else {

			}
			*/

			b, err := json.Marshal(data)
			if err != nil {
				return err
			}

			body := io.NopCloser(bytes.NewReader(b))
			response.Body = body
			contentLength := len(b)
			response.ContentLength = int64(contentLength)
			response.Header.Set("Content-Length", strconv.Itoa(contentLength))

			return nil
		}
	}

	return func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		rproxy.ServeHTTP(rw, r)
		return true, nil
	}
}

type ReverseProxyOption func(rp *reverseProxy)

// WithHostFromServiceUrl will not forward the original Host in the Headers
// but will overwrite this one while forwarding with the host of the target url
func WithHostFromServiceUrl() ReverseProxyOption {
	return func(rp *reverseProxy) {
		rp.useHostFromServiceUrl = true
	}
}

func ReverseProxy(serviceUrl string, options ...ReverseProxyOption) Action {
	u, err := url.Parse(serviceUrl)
	if err != nil {
		panic(fmt.Errorf("unable to parse url `%s` for reverse proxy: %w", serviceUrl, err))
	}

	rp := &reverseProxy{
		targetUrl: u,
	}

	for _, option := range options {
		option(rp)
	}

	return rp
}

type proxyResponseWriter struct {
	bytes.Buffer
}

func (p *proxyResponseWriter) Header() http.Header {
	//TODO implement me
	panic("implement me")
}

func (p *proxyResponseWriter) WriteHeader(statusCode int) {
	//TODO implement me
	panic("implement me")
}
