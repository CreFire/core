package core

import "reflect"

func Copy[T any](src T) T {
	return deepCopy(src, false).(T)
}

func deepCopy(src any, all bool) any {
	if src == nil {
		return nil
	}
	srcV := reflect.ValueOf(src)
	dstV := reflect.New(reflect.TypeOf(src)).Elem()
}
