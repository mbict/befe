package http

import (
	. "github.com/mbict/befe/expr"
	"io"
	"net/http"
	"net/url"
	"time"
)

type HttpClient interface {
	WithHeader(key string, value Valuer) HttpClient
	WithQueryParam(key string, value Valuer) HttpClient

	Get(path Valuer) Promise
	Post(path Valuer) Promise
	Put(path Valuer) Promise
	Patch(path Valuer) Promise
	Delete(path Valuer) Promise
	Head(path Valuer) Promise
}

type httpClient struct {
	baseUrl     string
	client      *http.Client
	queryParams map[string]Valuer
	headers     map[string]Valuer
}

func (c *httpClient) clone() *httpClient {

	q := make(map[string]Valuer, len(c.queryParams))
	for k, v := range c.queryParams {
		q[k] = v
	}

	h := make(map[string]Valuer, len(c.headers))
	for k, v := range c.headers {
		h[k] = v
	}

	return &httpClient{
		baseUrl:     c.baseUrl,
		client:      c.client,
		headers:     h,
		queryParams: q,
	}
}

func (c *httpClient) WithHeader(key string, value Valuer) HttpClient {
	c = c.clone()
	c.headers[key] = value
	return c
}

func (c *httpClient) WithQueryParam(key string, value Valuer) HttpClient {
	c = c.clone()
	c.queryParams[key] = value
	return c
}

func (c *httpClient) prepareRequest(r *http.Request, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	//generate the heades
	if len(c.headers) > 0 {
		req.Header = c.generateHeaders(r)
	}

	//generate the dynamic query params
	if len(c.queryParams) > 0 {
		req.URL.RawQuery = c.generateQuery(r, req.URL.Query())
	}
	return req, nil
}

func (c *httpClient) generateHeaders(r *http.Request) http.Header {
	headers := http.Header{}
	for header, valuer := range c.headers {
		headers.Add(header, ValueToString(valuer(r)))
	}
	return headers
}

func (c *httpClient) generateQuery(r *http.Request, queryParams url.Values) string {
	for queryName, valuer := range c.queryParams {
		queryParams.Set(queryName, ValueToString(valuer(r)))
	}
	return queryParams.Encode()
}

func (c *httpClient) Get(path Valuer) Promise {
	return c.createRequestHandler(func(r *http.Request) (*http.Response, error) {
		url := c.baseUrl + ValueToString(path(r))
		req, err := c.prepareRequest(r, "GET", url, nil)
		if err != nil {
			return nil, err
		}
		return c.client.Do(req)
	})
}

func (c *httpClient) Post(path Valuer) Promise {
	//TODO implement me
	panic("implement me post")
}

func (c *httpClient) Put(path Valuer) Promise {
	//TODO implement me
	panic("implement me put")
}

func (c *httpClient) Patch(path Valuer) Promise {
	//TODO implement me
	panic("implement me patch")
}

func (c *httpClient) Delete(path Valuer) Promise {
	//TODO implement me
	panic("implement me delete")
}

func (c *httpClient) Head(path Valuer) Promise {
	//TODO implement me
	panic("implement me head")
}

func (c *httpClient) createRequestHandler(call func(r *http.Request) (*http.Response, error)) Promise {
	return NewPromise(func(rw http.ResponseWriter, r *http.Request, success, failure Handler) (bool, error) {
		resp, err := call(r)
		if resp != nil {
			data, err := DecodeResponse(resp)
			if err != nil {
				return false, err
			}

			r = StoreResultBucket(r, "", data, resp, err)

			if err != nil {
				//error while reading
				return failure(rw, r)
			}
			//put result in context
			if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
				if success == nil {
					return true, nil
				}
				return success(rw, r)
			} else {
				if failure == nil {
					return true, nil
				}
				return failure(rw, r)
			}
		} else {
			//call failure, and store response error
			r = StoreResultBucket(r, "", nil, nil, err)

			if failure == nil {
				return true, nil
			}
			return failure(rw, r)
		}
	})
}

func Client(baseUrl string) HttpClient {
	return &httpClient{
		baseUrl: baseUrl,
		client: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   60 * time.Second,
		},
	}
}
