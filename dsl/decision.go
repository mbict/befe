package dsl

import (
	"context"
	"net/http"
)

type DecisionTree interface {
	Action

	//When creates a decision what action to perform if all conditions match
	When(...Condition) Decision

	//Default action is executed if all the decision tree fails to find an action
	Default(...Action) DecisionTree
}

type Decision interface {
	Action
	Condition
	When(...Condition) Decision
	Then(...Action) Decision
}

type decisionTree struct {
	decisionPaths  []Decision
	defaultActions Actions
}

func Decisions() DecisionTree {
	return &decisionTree{
		decisionPaths:  nil,
		defaultActions: nil,
	}
}

func (t *decisionTree) BuildHandler(ctx context.Context, next Handler) Handler {
	h := t.defaultActions.BuildHandler(ctx, next)
	for i := len(t.decisionPaths); i > 0; i-- {
		cond := t.decisionPaths[i-1].BuildCondition()
		act := t.decisionPaths[i-1].BuildHandler(ctx, next)
		nxt := h
		h = func(rw http.ResponseWriter, r *http.Request) {
			if cond(r) == true {
				act(rw, r)
				return
			}
			nxt(rw, r)
		}
	}
	return h
}

func (t *decisionTree) When(conditions ...Condition) Decision {
	d := newDecision().When(conditions...)
	t.decisionPaths = append(t.decisionPaths, d)
	return d
}

func (t *decisionTree) Default(actions ...Action) DecisionTree {
	t.defaultActions = append(t.defaultActions, actions...)
	return t
}

type decision struct {
	Actions
	Conditions
}

func newDecision() Decision {
	return &decision{}
}

func (d *decision) When(conditions ...Condition) Decision {
	d.Conditions = append(d.Conditions, conditions...)
	return d
}

func (d *decision) Then(actions ...Action) Decision {
	d.Actions = append(d.Actions, actions...)
	return d
}

type Breaker interface {
	Action

	When(...Condition) Breaker
	Default(...Action) Breaker
}

type breakerAction struct {
	mustBe       bool
	conditions   Conditions
	breakActions Actions
}

func (b *breakerAction) BuildHandler(ctx context.Context, next Handler) Handler {
	breakActions := b.breakActions.BuildHandler(ctx, func(rw http.ResponseWriter, r *http.Request) {})
	cond := b.conditions.BuildCondition()

	return func(rw http.ResponseWriter, r *http.Request) {
		if cond(r) == b.mustBe {
			next(rw, r)
			return
		}
		breakActions(rw, r)
	}
}

func (b *breakerAction) When(conditions ...Condition) Breaker {
	b.conditions = append(b.conditions, conditions...)
	return b
}

func (b *breakerAction) Default(actions ...Action) Breaker {
	b.breakActions = append(b.breakActions, actions...)
	return b
}

func StopWhen(conditions ...Condition) Breaker {
	return &breakerAction{
		mustBe:       false,
		conditions:   conditions,
		breakActions: nil,
	}
}

func ContinueWhen(conditions ...Condition) Breaker {
	return &breakerAction{
		mustBe:       true,
		conditions:   conditions,
		breakActions: nil,
	}
}
