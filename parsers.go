package env

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Required due to init cycle
func init() {
	parsers[reflect.Slice] = func(f reflect.Value, v string) error {
		sl := reflect.MakeSlice(f.Type(), 0, 0)
		if f.Type().Elem().Kind() == reflect.Uint8 {
			sl = reflect.ValueOf([]byte(v))
		} else if len(strings.TrimSpace(v)) != 0 {
			valCollection := strings.Split(v, ",")
			sl = reflect.MakeSlice(f.Type(), len(valCollection), len(valCollection))
			for i, v := range valCollection {
				err := setValue(sl.Index(i), v)
				if err != nil {
					return err
				}
			}
		}
		f.Set(sl)
		return nil
	}

	parsers[reflect.Map] = func(f reflect.Value, v string) error {
		mp := reflect.MakeMap(f.Type())
		if len(strings.TrimSpace(v)) != 0 {
			pairs := strings.Split(v, ",")
			for _, pair := range pairs {
				keyValues := strings.Split(pair, ":")
				if len(keyValues) != 2 {
					return fmt.Errorf("invalid map item: %q", pair)
				}
				k := reflect.New(f.Type().Key()).Elem()
				err := setValue(k, keyValues[0])
				if err != nil {
					return err
				}
				rv := reflect.New(f.Type().Elem()).Elem()
				err = setValue(rv, keyValues[1])
				if err != nil {
					return err
				}
				mp.SetMapIndex(k, rv)
			}
		}
		f.Set(mp)
		return nil
	}
}

type parser func(f reflect.Value, v string) error

var parsers = map[reflect.Kind]parser{
	reflect.String: func(f reflect.Value, v string) error {
		f.SetString(v)
		return nil
	},
	reflect.Bool: func(f reflect.Value, v string) error {
		val, err := strconv.ParseBool(v)
		if err != nil {
			return err
		}
		f.SetBool(val)
		return nil
	},
	reflect.Int: func(f reflect.Value, v string) error {
		return setInt(f, v, 32)
	},
	reflect.Int8: func(f reflect.Value, v string) error {
		return setInt(f, v, 8)
	},
	reflect.Int16: func(f reflect.Value, v string) error {
		return setInt(f, v, 16)
	},
	reflect.Int32: func(f reflect.Value, v string) error {
		return setInt(f, v, 32)
	},
	reflect.Int64: func(f reflect.Value, v string) error {
		if f.Type().PkgPath() == "time" && f.Type().Name() == "Duration" {
			d, err := time.ParseDuration(v)
			if err != nil {
				return err
			}
			f.SetInt(int64(d))
			return nil
		}
		return setInt(f, v, 64)
	},
	reflect.Uint: func(f reflect.Value, v string) error {
		return setUint(f, v, 32)
	},
	reflect.Uint8: func(f reflect.Value, v string) error {
		return setUint(f, v, 8)
	},
	reflect.Uint16: func(f reflect.Value, v string) error {
		return setUint(f, v, 16)
	},
	reflect.Uint32: func(f reflect.Value, v string) error {
		return setUint(f, v, 32)
	},
	reflect.Uint64: func(f reflect.Value, v string) error {
		return setUint(f, v, 64)
	},
	reflect.Float32: func(f reflect.Value, v string) error {
		return setFloat(f, v, 32)
	},
	reflect.Float64: func(f reflect.Value, v string) error {
		return setFloat(f, v, 64)
	},
}

func setInt(f reflect.Value, v string, bitSize int) error {
	val, err := strconv.ParseInt(v, 0, bitSize)
	if err != nil {
		return err
	}
	f.SetInt(val)
	return nil
}

func setUint(f reflect.Value, v string, bitSize int) error {
	val, err := strconv.ParseUint(v, 0, bitSize)
	if err != nil {
		return err
	}
	f.SetUint(val)
	return nil
}

func setFloat(f reflect.Value, v string, bitSize int) error {
	val, err := strconv.ParseFloat(v, bitSize)
	if err != nil {
		return err
	}
	f.SetFloat(val)
	return nil
}

func lookup(s string) (string, error) {
	val, ok := os.LookupEnv(s)
	if !ok {
		return val, fmt.Errorf("env: '%s' not found", s)
	}
	return val, nil
}

func AsString(s string) (string, error) {
	return lookup(s)
}

func AsBool(s string) (v bool, e error) {
	val, err := lookup(s)
	if err != nil {
		return v, err
	}
	v, e = strconv.ParseBool(val)
	if e != nil {
		return v, parseError(s, val, "bool", 0)
	}
	return v, nil
}

func AsInt(s string, bitSize int) (v int64, e error) {
	val, err := lookup(s)
	if err != nil {
		return v, err
	}
	v, e = strconv.ParseInt(val, 0, bitSize)
	if e != nil {
		return v, parseError(s, val, "int", bitSize)
	}
	return v, nil
}

func AsDuration(s string) (v time.Duration, e error) {
	val, err := lookup(s)
	if err != nil {
		return v, err
	}
	v, e = time.ParseDuration(val)
	if e != nil {
		return v, parseError(s, val, "duration", 0)
	}
	return v, nil
}

func AsFloat(s string, bitSize int) (v float64, e error) {
	val, err := lookup(s)
	if err != nil {
		return v, err
	}
	v, e = strconv.ParseFloat(val, bitSize)
	if e != nil {
		return v, parseError(s, val, "float", bitSize)
	}
	return v, nil
}

func AsUint(s string, bitSize int) (v uint64, e error) {
	val, err := lookup(s)
	if err != nil {
		return v, err
	}
	v, e = strconv.ParseUint(val, 0, bitSize)
	if e != nil {
		return v, parseError(s, val, "uint", bitSize)
	}
	return v, nil
}

func parseError(s string, v string, t string, b int) error {
	if b != 0 {
		t = fmt.Sprintf("%s[%d]", t, b)
	}
	return fmt.Errorf("env: unable to parse ['%s'='%s'] as %s", s, v, t)
}
