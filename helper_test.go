package env_test

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/halorium/env"
)

type runFunc func(t *testing.T) (interface{}, error)

func TestHelperFunctions(t *testing.T) {
	cases := []struct {
		name   string
		setEnv setEnv
		opts   env.Options
		run    runFunc
		err    error
		want   interface{}
	}{
		{
			name:   "AsString env not set",
			setEnv: func(t *testing.T) {},
			opts:   env.Options{},
			run: func(t *testing.T) (interface{}, error) {
				return env.AsString("ENV_VAR")
			},
			err:  fmt.Errorf("env: 'ENV_VAR' not found"),
			want: "string_val",
		},
		{
			name: "AsString valid env set",
			setEnv: func(t *testing.T) {
				t.Setenv("ENV_VAR", "string_val")
			},
			opts: env.Options{},
			run: func(t *testing.T) (interface{}, error) {
				return env.AsString("ENV_VAR")
			},
			err:  nil,
			want: "string_val",
		},
		{
			name:   "AsBool env not set",
			setEnv: func(t *testing.T) {},
			opts:   env.Options{},
			run: func(t *testing.T) (interface{}, error) {
				return env.AsBool("ENV_VAR")
			},
			err:  fmt.Errorf("env: 'ENV_VAR' not found"),
			want: true,
		},
		{
			name: "AsBool valid env set",
			setEnv: func(t *testing.T) {
				t.Setenv("ENV_VAR", "true")
			},
			opts: env.Options{},
			run: func(t *testing.T) (interface{}, error) {
				return env.AsBool("ENV_VAR")
			},
			err:  nil,
			want: true,
		},
		{
			name: "AsBool invalid env set",
			setEnv: func(t *testing.T) {
				t.Setenv("ENV_VAR", "invalid")
			},
			opts: env.Options{},
			run: func(t *testing.T) (interface{}, error) {
				return env.AsBool("ENV_VAR")
			},
			err:  fmt.Errorf("env: unable to parse ['ENV_VAR'='invalid'] as bool"),
			want: true,
		},
		{
			name:   "AsInt env not set",
			setEnv: func(t *testing.T) {},
			opts:   env.Options{},
			run: func(t *testing.T) (interface{}, error) {
				return env.AsInt("ENV_VAR", 64)
			},
			err:  fmt.Errorf("env: 'ENV_VAR' not found"),
			want: int64(1),
		},
		{
			name: "AsInt valid env set",
			setEnv: func(t *testing.T) {
				t.Setenv("ENV_VAR", "1")
			},
			opts: env.Options{},
			run: func(t *testing.T) (interface{}, error) {
				return env.AsInt("ENV_VAR", 64)
			},
			err:  nil,
			want: int64(1),
		},
		{
			name: "AsInt invalid env set",
			setEnv: func(t *testing.T) {
				t.Setenv("ENV_VAR", "invalid")
			},
			opts: env.Options{},
			run: func(t *testing.T) (interface{}, error) {
				return env.AsInt("ENV_VAR", 64)
			},
			err:  fmt.Errorf("env: unable to parse ['ENV_VAR'='invalid'] as int[64]"),
			want: int64(1),
		},
		{
			name:   "AsFloat env not set",
			setEnv: func(t *testing.T) {},
			opts:   env.Options{},
			run: func(t *testing.T) (interface{}, error) {
				return env.AsFloat("ENV_VAR", 64)
			},
			err:  fmt.Errorf("env: 'ENV_VAR' not found"),
			want: float64(1),
		},
		{
			name: "AsFloat valid env set",
			setEnv: func(t *testing.T) {
				t.Setenv("ENV_VAR", "1.2")
			},
			opts: env.Options{},
			run: func(t *testing.T) (interface{}, error) {
				return env.AsFloat("ENV_VAR", 64)
			},
			err:  nil,
			want: float64(1.2),
		},
		{
			name: "AsFloat invalid env set",
			setEnv: func(t *testing.T) {
				t.Setenv("ENV_VAR", "invalid")
			},
			opts: env.Options{},
			run: func(t *testing.T) (interface{}, error) {
				return env.AsFloat("ENV_VAR", 64)
			},
			err:  fmt.Errorf("env: unable to parse ['ENV_VAR'='invalid'] as float[64]"),
			want: float64(1.2),
		},
		{
			name:   "AsUint env not set",
			setEnv: func(t *testing.T) {},
			opts:   env.Options{},
			run: func(t *testing.T) (interface{}, error) {
				return env.AsUint("ENV_VAR", 64)
			},
			err:  fmt.Errorf("env: 'ENV_VAR' not found"),
			want: uint64(1),
		},
		{
			name: "AsUint valid env set",
			setEnv: func(t *testing.T) {
				t.Setenv("ENV_VAR", "1")
			},
			opts: env.Options{},
			run: func(t *testing.T) (interface{}, error) {
				return env.AsUint("ENV_VAR", 64)
			},
			err:  nil,
			want: uint64(1),
		},
		{
			name: "AsUint invalid env set",
			setEnv: func(t *testing.T) {
				t.Setenv("ENV_VAR", "invalid")
			},
			opts: env.Options{},
			run: func(t *testing.T) (interface{}, error) {
				return env.AsUint("ENV_VAR", 64)
			},
			err:  fmt.Errorf("env: unable to parse ['ENV_VAR'='invalid'] as uint[64]"),
			want: uint64(1),
		},
		{
			name:   "AsDuration env not set",
			setEnv: func(t *testing.T) {},
			opts:   env.Options{},
			run: func(t *testing.T) (interface{}, error) {
				return env.AsDuration("ENV_VAR")
			},
			err:  fmt.Errorf("env: 'ENV_VAR' not found"),
			want: time.Duration(5 * time.Second),
		},
		{
			name: "AsDuration valid env set",
			setEnv: func(t *testing.T) {
				t.Setenv("ENV_VAR", "5s")
			},
			opts: env.Options{},
			run: func(t *testing.T) (interface{}, error) {
				return env.AsDuration("ENV_VAR")
			},
			err:  nil,
			want: time.Duration(5 * time.Second),
		},
		{
			name: "AsDuration invalid env set",
			setEnv: func(t *testing.T) {
				t.Setenv("ENV_VAR", "invalid")
			},
			opts: env.Options{},
			run: func(t *testing.T) (interface{}, error) {
				return env.AsDuration("ENV_VAR")
			},
			err:  fmt.Errorf("env: unable to parse ['ENV_VAR'='invalid'] as duration"),
			want: time.Duration(5 * time.Second),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.setEnv(t)
			got, err := c.run(t)

			if err != nil && c.err != nil {
				if err.Error() != c.err.Error() {
					t.Errorf("\nwant:'%#v'\ngot:'%#v'\n", c.err.Error(), err.Error())
				}
			} else if err != c.err {
				t.Errorf("\nwant:'%#v'\ngot:'%#v'\n", c.err, err)
			}

			if err == nil && c.err == nil {
				if !reflect.DeepEqual(c.want, got) {
					t.Errorf("\nwant:'%#v'\ngot:'%#v'\n", c.want, got)
				}
			}

		})
	}
}
