package expr

import (
	"context"
	"errors"
	"net/http"
)

type Transformer interface {
	Action
	Transform(interface{}) interface{}
}

type TransformerFunc func(interface{}) interface{}

func (t TransformerFunc) BuildHandler(ctx context.Context, next Handler) Handler {
	return ActionFunc(func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		if res := GetResultBucket(r.Context()); res != nil {
			res.Data = t(res.Data)
			return true, nil
		}
		return false, errors.New("cannot perform transformation, there is no result body")
	}).BuildHandler(ctx, next)
}

func (t TransformerFunc) Transform(i interface{}) interface{} {
	return t(i)
}
