package expr

import (
	"context"
	"net/http"
)

type DecisionTree interface {
	Action
	When(...Condition) DecisionCondition
	ElseCondition
}

type DecisionCondition interface {
	Then(...Action) DecisionTree
}

type decisionTree struct {
	decisions   []*decisionCondition
	elseActions Actions
}

func (d *decisionTree) BuildHandler(ctx context.Context, next Handler) Handler {
	h := d.elseActions.BuildHandler(ctx, next) // <-- we are nesting next here on purpose

	if h == nil {
		h = EmptyHandler
	}

	for i := len(d.decisions); i > 0; i-- {
		cond := d.decisions[i-1].BuildCondition(ctx)
		act := d.decisions[i-1].BuildHandler(ctx, next) // <-- we are nesting next here on purpose
		nxt := h
		h = func(rw http.ResponseWriter, r *http.Request) (bool, error) {
			if cond(r) == true {
				return act(rw, r)
			}
			return nxt(rw, r)
		}
	}
	return h
}

func (d *decisionTree) When(condition ...Condition) DecisionCondition {
	swd := &decisionCondition{
		decisionTree: d,
		Conditions:   Conditions(condition),
		Actions:      Actions{},
	}
	d.decisions = append(d.decisions, swd)
	return swd
}

func (d *decisionTree) Else(action ...Action) Action {
	d.elseActions = append(d.elseActions, action...)
	return d
}

type decisionCondition struct {
	Conditions
	Actions
	*decisionTree
}

func (s *decisionCondition) BuildHandler(ctx context.Context, next Handler) Handler {
	return s.Actions.BuildHandler(ctx, next)
}

func (s *decisionCondition) Then(action ...Action) DecisionTree {
	s.Actions = append(s.Actions, action...)
	return s.decisionTree
}

func NewDecisionTree() DecisionTree {
	return &decisionTree{
		decisions:   []*decisionCondition{},
		elseActions: Actions{},
	}
}
