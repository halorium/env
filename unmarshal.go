package env

import (
	"errors"
	"fmt"
	"os"
	"reflect"
)

var ErrInvalidType = errors.New("must be a pointer to a non-nil struct")

func Unmarshal(obj interface{}, options ...Options) error {
	// get default options and merge in any overrides
	opts := getOptions(options...)

	rv := reflect.ValueOf(obj)
	if rv.Kind() != reflect.Ptr {
		return ErrInvalidType
	}

	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return ErrInvalidType
	}

	// recurse the struct and set env fields
	return parseStruct(obj, opts)
}

func parseStruct(obj interface{}, opts Options) error {
	rv := reflect.ValueOf(obj)
	rt := rv.Type()

	if rt.Kind() == reflect.Ptr {
		rv = rv.Elem()
		rt = rv.Type()
	}

	if rt.Kind() != reflect.Struct {
		return fmt.Errorf("object must be a struct")
	}

	// iterate over struct fields
	for i := 0; i < rv.NumField(); i++ {
		rf := rv.Field(i)
		rsf := rv.Type().Field(i)

		// ignore non exported fields
		if !rf.CanSet() {
			continue
		}

		// if pointer to struct or nil struct (instantiate it)
		if rf.Kind() == reflect.Ptr && rf.Type().Elem().Kind() == reflect.Struct {
			if rf.IsNil() {
				// nil pointer to struct: create a zero instance
				rf.Set(reflect.New(rf.Type().Elem()))
			}
			rf = rf.Elem()
		}

		// if struct we need to recurse (unless implements Unmarshaler)
		if rf.Kind() == reflect.Struct && asUnmarshaler(rf) == nil {
			rfi := rf.Addr().Interface()
			err := parseStruct(rfi, opts)
			if err != nil {
				return err
			}
			continue
		}

		// ignore fields without a tag or explicitly ignored
		if rsf.Tag.Get(opts.Tag) == "-" || rsf.Tag.Get(opts.Tag) == "" {
			continue
		}

		tag := rsf.Tag.Get(opts.Tag)

		val, ok := os.LookupEnv(tag)
		if !ok {
			if opts.Required {
				return fmt.Errorf("'%s' is required", tag)
			}
			// skip it
			continue
		}

		// now we can parse
		err := setValue(rf, val)
		if err != nil {
			return err
		}
	}

	return nil
}

func setValue(rf reflect.Value, val string) error {
	// check for custom UnmarshalENV function
	if f := asUnmarshaler(rf); f != nil {
		return f.UnmarshalENV(val)
	}

	// instantiate field for pointer
	if rf.Type().Kind() == reflect.Ptr {
		if rf.IsNil() {
			rf.Set(reflect.New(rf.Type().Elem()))
		}
		rf = rf.Elem()
	}

	setter, ok := parsers[rf.Type().Kind()]
	if !ok {
		return fmt.Errorf("field:'%s' is an unsupported type:'%s'", rf, rf.Type().String())
	}
	return setter(rf, val)
}
