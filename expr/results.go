package expr

import (
	"context"
	"net/http"
)

var resultBucketContextKey int

type ResultBucket struct {
	Name     string
	Data     interface{}
	Error    error
	Response *http.Response

	previousBucket *ResultBucket
}

func (b *ResultBucket) PreviousBucket() *ResultBucket {
	return b.previousBucket
}

func GetResultBucket(ctx context.Context) *ResultBucket {
	if r, ok := ctx.Value(&resultBucketContextKey).(*ResultBucket); ok {
		return r
	}
	return nil
}

func StoreResultBucket(r *http.Request, name string, data interface{}, response *http.Response, err error) *http.Request {
	bucket := &ResultBucket{
		Name:           name,
		Data:           data,
		Response:       response,
		Error:          err,
		previousBucket: GetResultBucket(r.Context()),
	}
	return r.WithContext(context.WithValue(r.Context(), &resultBucketContextKey, bucket))
}
