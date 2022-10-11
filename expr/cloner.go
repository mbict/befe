package expr

type Cloner[T any] interface {
	Clone() T
}
