package env

import "reflect"

type Unmarshaler interface {
	UnmarshalENV(value string) error
}

func asUnmarshaler(rv reflect.Value) Unmarshaler {
	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
	} else if rv.CanAddr() {
		rv = rv.Addr()
	}
	if rv.CanInterface() {
		u, ok := rv.Interface().(Unmarshaler)
		if !ok {
			return nil
		}
		return u
	}
	return nil
}
