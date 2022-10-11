package expr

import (
	"context"
	"net/http"
)

type Promise interface {
	Action

	OnSuccess(...Action) Promise
	OnFailure(...Action) Promise
}

func NewPromise(handler func(rw http.ResponseWriter, r *http.Request, success, failure Handler) (bool, error)) Promise {
	return &promise{
		handler: handler,
		success: Actions{},
		failure: Actions{},
	}
}

type promise struct {
	handler func(rw http.ResponseWriter, r *http.Request, success, failure Handler) (bool, error)
	success Actions
	failure Actions
}

func (c *promise) Clone() Promise {
	return &promise{
		handler: c.handler,
		success: c.success,
		failure: c.failure,
	}
}

func (c *promise) BuildHandler(ctx context.Context, next Handler) Handler {
	successHandler := c.success.BuildHandler(ctx, nil)
	failureHandler := c.failure.BuildHandler(ctx, nil)

	return func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		if cont, err := c.handler(rw, r, successHandler, failureHandler); err != nil || cont == false {
			return cont, err
		}

		if next == nil {
			return true, nil
		}
		return next(rw, r)
	}
}

func (c *promise) OnSuccess(action ...Action) Promise {
	c.success = append(c.success, action...)
	return c
}

func (c *promise) OnFailure(action ...Action) Promise {
	c.failure = append(c.failure, action...)
	return c
}
