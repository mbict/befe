package http

import (
	"context"
	. "github.com/mbict/befe/dsl"
	. "github.com/mbict/befe/expr"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

var paginatorOffsetContextKey int
var paginatorCursorContextKey int

type offsetPaginator struct {
	sizePerRequest    int
	maxResultsFetcher Valuer
	response          Promise
}

func (o *offsetPaginator) BuildHandler(ctx context.Context, next Handler) Handler {
	failure := ActionFunc(func(writer http.ResponseWriter, r *http.Request) (bool, error) {
		//on failure we always stop
		*(r.Context().Value(&paginatorOffsetContextKey).(*int)) = -1
		return true, nil
	})

	success := ActionFunc(func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		maxResults, err := strconv.Atoi(ValueToString(o.maxResultsFetcher(r)))
		if err != nil {
			return false, errors.WithMessage(err, "cannot determine max results")
		}

		//determine if the next  request will go over the max results
		offset := r.Context().Value(&paginatorOffsetContextKey).(*int)
		if *offset+o.sizePerRequest > maxResults {
			*offset = -1
		}
		return true, nil
	})

	h := o.response.
		OnSuccess(success).
		OnFailure(failure).
		BuildHandler(ctx, nil)

	return func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		offset := 0
		r = r.WithContext(context.WithValue(r.Context(), &paginatorOffsetContextKey, &offset))
		for offset != -1 {
			if c, err := h(rw, r); err != nil || c == false {
				return c, err
			}
			if offset >= 0 {
				offset += o.sizePerRequest
			}
		}

		if next == nil {
			return true, nil
		}
		return next(rw, r)
	}

}

//func (o *offsetPaginator) OnSuccess(action ...Action) Promise {
//	o.success = append(o.success, action...)
//	return o
//}
//
//func (o *offsetPaginator) OnFailure(action ...Action) Promise {
//	o.failure = append(o.failure, action...)
//	return o
//}

// OffsetPaginatedResults is a helper that will fetch all pages from a paginated api endpoint and combine the results into one set.
func OffsetPaginatedResults(sizePerRequest int, maxResultsFetcher Valuer, response Promise) Action {
	//try to clone the promise
	if cloner, ok := response.(Cloner[Promise]); ok {
		response = cloner.Clone()
	}

	return &offsetPaginator{
		sizePerRequest:    sizePerRequest,
		maxResultsFetcher: maxResultsFetcher,
		response:          response,
	}
}

// CursorPaginatedREsults is a helper that will fetch all pages from a cursor paginated api endpoint and combine the results into one set.
func CursorPaginatedResults(nextCursorFetcher Valuer, response Promise) Promise {
	panic("please implement me")
}

func ParamFromPaginator(paramName, paginatorParam string) Param {
	return WithParam(paramName, ValueFromPaginator(paginatorParam))
}

func ValueFromPaginator(paginatorParam string) Valuer {
	switch paginatorParam {
	case "offset":
		return func(r *http.Request) interface{} {
			return *(r.Context().Value(&paginatorOffsetContextKey)).(*int)
		}

	case "cursor":
		return func(r *http.Request) interface{} {
			return r.Context().Value(&paginatorCursorContextKey)
		}
	}
	panic("invalid param name for paginator")
}
