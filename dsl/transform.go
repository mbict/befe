package dsl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type TransformFunc func(*http.Response) error

func (t TransformFunc) Build() TransformFunc {
	return t
}

type Transformers []Transformer

func (t Transformers) Build() TransformFunc {

	//if there is only one in the stack we do not wrap it
	if len(t) == 1 {
		return t[0].Build()
	}

	transformStack := []TransformFunc{}
	for _, transformer := range t {
		transformer.Build()
	}

	return func(response *http.Response) error {
		for _, transform := range transformStack {
			if err := transform(response); err != nil {
				return err
			}
		}
		return nil
	}
}

type Transformer interface {
	Build() TransformFunc
}

type JsonTransformer func(message map[string]interface{}) error

type ResponseModifier struct {
	rw http.ResponseWriter
}

func (r *ResponseModifier) Header() http.Header {
	panic("implement me")
}

func (r *ResponseModifier) Write(bytes []byte) (int, error) {
	panic("implement me")
}

func (r *ResponseModifier) WriteHeader(statusCode int) {
	panic("implement me")
}

//--- helpers
func Json(transformers ...JsonTransformer) Transformer {
	return TransformFunc(func(response *http.Response) error {
		//decode body
		decoder := json.NewDecoder(response.Body)
		document := make(map[string]interface{})
		if err := decoder.Decode(&document); err != nil {
			return err
		}

		//do transformations
		for _, transform := range transformers {
			if err := transform(document); err != nil {
				return err
			}
		}

		//transform back
		buf := bytes.NewBuffer(nil)
		encoder := json.NewEncoder(buf)
		if err := encoder.Encode(document); err != nil {
			return err
		}

		//close the old body
		response.Body.Close()

		//create a new body from the result
		response.Body = ioutil.NopCloser(buf)
		size := buf.Len()
		response.ContentLength = int64(size)
		response.Header["Content-Length"] = []string{strconv.Itoa(size)}

		return nil
	})
}

type patternTree map[string]patternTree

func copyMap(src map[string]patternTree, dst map[string]patternTree) {
	for k, v := range src {
		if v == nil {
			dst[k] = nil
		}

		if _, ok := dst[k]; !ok {
			dst[k] = make(patternTree)
		}

		copyMap(v, dst[k])
	}
}

func createWildcardTree(tree patternTree) patternTree {

	wildcardTree := patternTree{}

	//has wildcard
	wildcardChilds, hasWildcard := tree["*"]
	if !hasWildcard {
		wildcardChilds = make(patternTree)
	}

	for elemName, childs := range tree {
		//we skip the wildcards
		if elemName == "*" {
			continue
		}

		wildcardTree[elemName] = make(patternTree)
		copyMap(wildcardChilds, wildcardTree[elemName])
		copyMap(childs, wildcardTree[elemName])

		//do it recursive
		wildcardTree[elemName] = createWildcardTree(wildcardTree[elemName])

	}

	if hasWildcard {
		wildcardTree["*"] = make(patternTree)
		copyMap(wildcardChilds, wildcardTree["*"])

		//do it recursive
		wildcardTree["*"] = createWildcardTree(wildcardTree["*"])

	}

	return wildcardTree
}

func Filter(fields ...string) JsonTransformer {

	var traverseTree func(data interface{}, tree patternTree) interface{}
	traverseTree = func(data interface{}, tree patternTree) interface{} {

		if tree == nil || len(tree) == 0 {
			return data
		}

		if dataSlice, ok := data.([]interface{}); ok {

			for index := len(dataSlice) - 1; index >= 0; index-- {
				childs, hasPattern := tree[strconv.Itoa(index)]
				if !hasPattern {
					childs, hasPattern = tree[`*`]
				}

				if hasPattern {
					dataSlice[index] = traverseTree(dataSlice[index], childs)
				} else {
					//remove slice, maintaining order
					dataSlice = append(dataSlice[:index], dataSlice[index+1:]...)
				}
			}
			data = dataSlice
		} else if dataMap, ok := data.(map[string]interface{}); ok {
			for key, _ := range dataMap {
				childs, hasPattern := tree[key]
				if !hasPattern {
					childs, hasPattern = tree[`*`]
				}

				if hasPattern {
					dataMap[key] = traverseTree(dataMap[key], childs)

				} else {
					//not found we remove
					delete(dataMap, key)
				}
			}
		} else {
			fmt.Println("not a thing i can use")
		}

		return data
	}

	//extract the raw pattern path
	tree := make(patternTree)
	for _, pattern := range fields {
		currentMap := tree
		pathNames := strings.Split(pattern, `.`)
		for _, name := range pathNames {
			if currentMap[name] == nil {
				currentMap[name] = make(patternTree)
			}
			currentMap = currentMap[name]
		}
	}
	//fill in the wildcards
	tree = createWildcardTree(tree)

	return func(data map[string]interface{}) error {
		traverseTree(data, tree)
		return nil
	}
}

//--- from the action stack add a transformer to modify the response
func TransformResponse(transformers ...Transformer) Action {
	transformer := Transformers(transformers).Build()
	a := transformResponse(transformer)
	return &a
}

var transformContextKey int = 1

type transformResponse TransformFunc

func (a *transformResponse) BuildHandler(_ context.Context, next Handler) Handler {
	return func(rw http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), &transformContextKey, TransformFunc(*a))
		r = r.WithContext(ctx)
		if next != nil {
			next(rw, r)
		}
	}
}
