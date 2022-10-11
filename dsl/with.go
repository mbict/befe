package dsl

import . "github.com/mbict/befe/expr"

// With is used to chain multiple actions into one executable chain
func With(actions ...Action) Action {
	return Actions(actions)
}
