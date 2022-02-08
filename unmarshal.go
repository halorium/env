package env

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidType = errors.New("must be a pointer to a non-nil struct")

const defaultTag = "env"

var defaultOptions = Options{
	Tag:      defaultTag,
	Required: false,
}

type Options struct {
	Tag string // default "env"
	Required bool // default false
}

type Unmarshaler interface {
	UnmarshalEnv(value string) error
}

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

	// recurse the struct and get env var data
	envVars := make([]envVar,0)
	envVars, err := parseStruct(obj, opts)
	if err != nil {
		return err
	}

	for _, envVar := range envVars {
		var val string
		var ok bool
		val, ok = os.LookupEnv(envVar.TagKey)
		if !ok {
			if opts.Required {
				return fmt.Errorf("'%s' is required", envVar.TagKey)
			}
			// skip it
			continue
		}

		err = setValue(envVar.Field, val)
		if err != nil {
			return err
		}
	}

	return nil
}

func getOptions(opts ...Options) Options {
	o := defaultOptions
	for _, opt := range opts {
		if opt.Tag != "" && opt.Tag != defaultTag {
			o.Tag = opt.Tag
		}
		if opt.Required {
			o.Required = opt.Required
		}
	}
	return o
}

type envVar struct {
	FieldName  string
	TagKey   string
	Field reflect.Value
	Tags  reflect.StructTag
}

func parseStruct(obj interface{}, opts Options) ([]envVar,error) {
	rv := reflect.ValueOf(obj)
	rt := rv.Type()

	if rt.Kind() == reflect.Ptr {
		rv = rv.Elem()
		rt = rv.Type()
	}

	if rt.Kind() != reflect.Struct {
		return nil, fmt.Errorf("object must be a struct")
	}

	envVars := make([]envVar,0)

	// iterate over struct fields
	for i := 0; i < rv.NumField(); i++ {
		rf := rv.Field(i)
		rsf := rv.Type().Field(i)

		// ignore non exported fields
		if !rf.CanSet() {
			continue
		}

		// if pointer to struct or nil struct (instantiate it)
		if rf.Kind() == reflect.Ptr && rf.Type().Elem().Kind() == reflect.Struct{
			if rf.IsNil() {
				// nil pointer to struct: create a zero instance
				rf.Set(reflect.New(rf.Type().Elem()))
			}

			rf = rf.Elem()
		}

		// if struct we need to recurse
		if rf.Kind() == reflect.Struct {
			rfi := rf.Addr().Interface()
			envVrs, err := parseStruct(rfi, opts)
			if err != nil {
				return nil, err
			}
			envVars = append(envVars, envVrs...)
			continue
		}

		// non struct objects below

		// ignore fields without a tag or explicitly ignored
		if rsf.Tag.Get(opts.Tag) == "-" || rsf.Tag.Get(opts.Tag) == ""{
			continue
		}

		varInfo := envVar{
			FieldName:  rsf.Name,
			TagKey: rsf.Tag.Get(opts.Tag),
			Field: rf,
			Tags:  rsf.Tag,
		}

		envVars = append(envVars, varInfo)
	}

	return envVars, nil
}

func setValue(rv reflect.Value, val string) error {
	rt := rv.Type()

	// check for custom UnmarshalEnv function
	if rv.CanInterface() {
		u, ok := rv.Interface().(Unmarshaler)
		if ok {
			return u.UnmarshalEnv(val)
		}
	}

	// instantiate field for pointer
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		if rv.IsNil() {
			rv.Set(reflect.New(rt))
		}
		rv = rv.Elem()
	}

	switch rt.Kind() {
	case reflect.String:
		rv.SetString(val)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if rv.Kind() == reflect.Int64 && rt.PkgPath() == "time" && rt.Name() == "Duration" {
			d, err := time.ParseDuration(val)
			if err != nil {
				return err
			}
			rv.SetInt(int64(d))
		} else {
			v, err := strconv.ParseInt(val, 0, rt.Bits())
			if err != nil {
				return err
			}
			rv.SetInt(v)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v, err := strconv.ParseUint(val, 0, rt.Bits())
		if err != nil {
			return err
		}
		rv.SetUint(v)

	case reflect.Bool:
		v, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		rv.SetBool(v)

	case reflect.Float32, reflect.Float64:
		v, err := strconv.ParseFloat(val, rt.Bits())
		if err != nil {
			return err
		}
		rv.SetFloat(v)

	case reflect.Slice:
		sl := reflect.MakeSlice(rt, 0, 0)
		if rt.Elem().Kind() == reflect.Uint8 {
			sl = reflect.ValueOf([]byte(val))
		} else if len(strings.TrimSpace(val)) != 0 {
			valCollection := strings.Split(val, ",")
			sl = reflect.MakeSlice(rt, len(valCollection), len(valCollection))
			for i, v := range valCollection {
				err := setValue(sl.Index(i), v)
				if err != nil {
					return err
				}
			}
		}
		rv.Set(sl)

	case reflect.Map:
		mp := reflect.MakeMap(rt)
		if len(strings.TrimSpace(val)) != 0 {
			pairs := strings.Split(val, ",")
			for _, pair := range pairs {
				keyValues := strings.Split(pair, ":")
				if len(keyValues) != 2 {
					return fmt.Errorf("invalid map item: %q", pair)
				}
				k := reflect.New(rt.Key()).Elem()
				err := setValue(k, keyValues[0])
				if err != nil {
					return err
				}
				v := reflect.New(rt.Elem()).Elem()
				err = setValue(v, keyValues[1])
				if err != nil {
					return err
				}
				mp.SetMapIndex(k, v)
			}
		}
		rv.Set(mp)

	default:
		return fmt.Errorf("field:'%s' is an unsupported type:'%s'", rv, rv.Type().String())
	}

	return nil
}

func lookup(k string) (string, error){
	val, ok := os.LookupEnv(k)
	if !ok {
		return val, fmt.Errorf("'%s' not found",k)
	}
	return val, nil
}

func AsString(k string) (string, error) {
	return lookup(k)
}

func AsBool(k string) (b bool, err error) {
	val, err := lookup(k)
	if err != nil {
		return b, err
	}
	return strconv.ParseBool(val)
}

func AsInt(k string, bitSize int) (i int64, err error) {
	val, err := lookup(k)
	if err != nil {
		return i, err
	}
	return strconv.ParseInt(val, 0, bitSize)
}

func AsDuration(k string) (d time.Duration, err error) {
	val, err := lookup(k)
	if err != nil {
		return d, err
	}
	return time.ParseDuration(val)
}

func AsFloat(k string, bitSize int) (f float64, err error) {
	val, err := lookup(k)
	if err != nil {
		return f, err
	}
	return strconv.ParseFloat(val, bitSize)
}

func AsUint(k string, bitSize int) (u uint64, err error) {
	val, err := lookup(k)
	if err != nil {
		return u, err
	}
	return strconv.ParseUint(val, 0, bitSize)
}