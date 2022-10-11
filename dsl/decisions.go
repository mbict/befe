package dsl

import (
	"github.com/mbict/befe/expr"
)

func Decision() expr.DecisionTree {
	return expr.NewDecisionTree()
}
