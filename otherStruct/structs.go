package otherStruct

type Stack[T comparable] interface {
	Push(value T)
	Pop() (value T, ok bool)
	Peek() (value T, ok bool)

	//containers.Container[T]
	Empty() bool
	Size() int
	Clear()
	Values() []interface{}
}

type Container[T any] interface {
	Empty() bool
	Size() int
	Clear()
	Values() []T
}
