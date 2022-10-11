package dsl

import (
	"context"
	"errors"
	"fmt"
	. "github.com/mbict/befe/expr"
	"github.com/ohler55/ojg/jp"
	"net/http"
	"strconv"
	"strings"
)

type transform []Transformer

func (t transform) Transform(i interface{}) interface{} {
	for _, transformer := range t {
		i = transformer.Transform(i)
	}
	return i
}

func (t transform) BuildHandler(ctx context.Context, next Handler) Handler {
	return ActionFunc(func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		if res := GetResultBucket(r.Context()); res != nil {
			res.Data = t.Transform(res.Data)
			return true, nil
		}
		return false, errors.New("cannot perform transformation, there is no result body")
	}).BuildHandler(ctx, next)
}

func Transform(tf ...Transformer) Action {
	return transform(tf)
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

// IncludePath is a transformer that will only keep the fields/values that match the path(s)
// Patterns:
//
//	data.*.id
//	data.arrayfield.2.*
//	data.name
func IncludePath(pathPatterns ...string) Transformer {

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
	for _, pattern := range pathPatterns {
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

	return TransformerFunc(func(in interface{}) interface{} {
		return traverseTree(in, tree)
	})
}

// ExcludePath is a transformer that will remove all the fields/values the match the path(s)
func ExcludePath(path ...string) Transformer {
	panic("")
}

// JsonPath is a transformer that will replace the buffered response with the result of the jsonPath
func JsonPath(pattern string) Transformer {
	panic("")
}

// JsonPath is a transformer that will replace the buffered response with the first result of the jsonPath
func JsonPathFirst(pattern string) Transformer {
	jq, err := jp.ParseString(pattern)
	if err != nil {
		panic("cannot compile JsonPath pattern `" + pattern + "`: " + err.Error())
	}

	return TransformerFunc(jq.First)
}
