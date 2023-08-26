package core

import "reflect"

func Copy[T any](src T) T {
	return deepCopy(src, false).(T)
}

func deepCopy(src any, all bool) (dst any) {
	if src == nil {
		return nil
	}
	srcValue := reflect.ValueOf(src)
	dstValue := reflect.New(srcValue.Type()).Elem()
	parent := make(map[uintptr]struct{})
	copyValue(dstValue, srcValue, all, parent)
	return dstValue.Interface()

}
func copyValue(dst reflect.Value, src reflect.Value, all bool, parent map[uintptr]struct{}) {
	if !dst.CanSet() {
		return
	}

	switch src.Kind() {
	case reflect.Struct:
		for i := 0; i < src.NumField(); i++ {
			copyValue(dst.Field(i), src.Field(i), all, parent)
		}
	case reflect.Array:
		for i := 0; i < src.Len(); i++ {
			copyValue(dst.Index(i), src.Index(i), all, parent)
		}
	case reflect.Slice:
		slice := reflect.MakeSlice(src.Type(), src.Len(), src.Cap())
		dst.Set(slice)
		for i := 0; i < src.Len(); i++ {
			copyValue(dst.Index(i), src.Index(i), all, parent)
		}
	case reflect.Map:
		newMap := reflect.MakeMap(src.Type())
		dst.Set(newMap)
		for _, key := range src.MapKeys() {
			value := src.MapIndex(key)
			newValue := reflect.New(value.Type()).Elem()
			copyValue(newValue, src.MapIndex(key), all, parent)
			dst.SetMapIndex(key, newValue)
		}
	case reflect.Func:
		if src.IsNil() {
			return
		}
		if !checkCanSet(dst, all) {
			return
		}
	case reflect.Ptr:
		if src.IsNil() {
			return
		}
		ptrValue := reflect.New(src.Elem().Type())
		copyValue(ptrValue.Elem(), src.Elem(), all, parent)
		dst.Set(ptrValue)
	default:
		dst.Set(src)
	}
}
func checkCanSet(dst reflect.Value, all bool) bool {
	if all {
		return true
	}
	return dst.CanSet()
}
