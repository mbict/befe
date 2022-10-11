package expr

type ConditionMiddleware interface {
	Action

	Must(...Condition) ConditionMiddleware
	OnFailure(...Action) ConditionMiddleware
}
