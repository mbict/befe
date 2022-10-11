package dsl

import (
	"context"
	"fmt"
	"github.com/mbict/befe/expr"
	"github.com/ohler55/ojg/jp"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func GetResult() expr.Valuer {
	return func(r *http.Request) interface{} {
		b := expr.GetResultBucket(r.Context())
		if b != nil {
			return b.Data
		}
		return nil
	}
}

type Merger interface {
	expr.Action

	Source(sourcer expr.Valuer) Merger
	Target(path string) Merger

	// Matcher gets a jsonPath pattern to find and match nodes
	Matcher(sourcePattern string, targetPattern string) Merger
}

func ResultsetMerger() Merger {
	return &resultsetMerger{}
}

type resultsetMerger struct {
	targetPath string
	sourcer    expr.Valuer

	//matcher jsonPath patterns
	sourceMatchPattern jp.Expr
	targetMatchPattern jp.Expr
}

func (m *resultsetMerger) BuildHandler(ctx context.Context, next expr.Handler) expr.Handler {
	return func(rw http.ResponseWriter, r *http.Request) (bool, error) {
		source := expr.GetResultBucket(r.Context())
		target := source.PreviousBucket()

		sourceData := m.sourcer(r)

		dataset := map[interface{}][]interface{}{}
		for _, data := range sourceData.([]interface{}) {
			lookupKey := m.sourceMatchPattern.First(data)

			dataset[lookupKey] = append(dataset[lookupKey], data)
		}

		targetField := filepath.Ext(m.targetPath)
		targetPath := strings.TrimSuffix(m.targetPath, targetField)
		targetField = strings.Trim(targetField, ` .`)

		jsonPatternWalker(targetPath, target.Data, func(data interface{}) {
			if dm, ok := data.(map[string]interface{}); ok {
				lookupKey := m.targetMatchPattern.First(dm)
				if mergeData, ok := dataset[lookupKey]; ok {
					mergeInto, ok := dm[targetField].([]interface{})
					if !ok {
						dm[targetField] = mergeData
					} else {
						dm[targetField] = append(mergeInto, mergeData...)
					}
				}

			} else {
				fmt.Println("cannot work with this type of node in json iterator")
			}
		})
		if next == nil {
			return true, nil
		}
		return next(rw, r)
	}
}

func (m *resultsetMerger) Source(sourcer expr.Valuer) Merger {
	m.sourcer = sourcer
	return m
}

func (m *resultsetMerger) Target(path string) Merger {
	m.targetPath = path
	return m
}

func (m *resultsetMerger) Matcher(sourcePattern string, targetPattern string) Merger {
	var err error
	if m.sourceMatchPattern, err = jp.ParseString(sourcePattern); err != nil {
		panic(err)
	}
	if m.targetMatchPattern, err = jp.ParseString(targetPattern); err != nil {
		panic(err)
	}
	return m
}

func jsonPatternWalker(path string, data interface{}, visitorFunc func(data interface{})) {

	var traverseTree func(data interface{}, path []string) interface{}
	traverseTree = func(data interface{}, path []string) interface{} {

		if path == nil || len(path) == 0 {
			//fmt.Println("arrived")
			visitorFunc(data)
			return data
		}

		currentPath := path[0]

		if dataSlice, ok := data.([]interface{}); ok {
			for index := len(dataSlice) - 1; index >= 0; index-- {
				if currentPath == strconv.Itoa(index) { //match on index
					//fmt.Println("here slice index", currentPath)
					traverseTree(dataSlice[index], path[1:])
					break
				} else if currentPath == "*" {
					//fmt.Println("here slice wildcard", currentPath)
					traverseTree(dataSlice[index], path[1:])
				}
			}
			data = dataSlice
		} else if dataMap, ok := data.(map[string]interface{}); ok {
			if currentPath == "*" { //wildcard all nodes
				for _, subData := range dataMap {
					//fmt.Println("here map wildcard", currentPath)
					traverseTree(subData, path[1:])
				}
			} else if subData, ok := dataMap[currentPath]; ok { //single named node
				//fmt.Println("here map index", currentPath)
				traverseTree(subData, path[1:])
			}
		} else {
			//fmt.Println("not a thing i can use", data)
		}

		return data
	}

	traverseTree(data, strings.Split(path, `.`))
}
