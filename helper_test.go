package env_test

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
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

			if err == nil {
				if diff := cmp.Diff(c.want, got); diff != "" {
					t.Errorf("(-want +got):\n%s", diff)
				}
				// if !reflect.DeepEqual(c.want, got) {
				// 	t.Errorf("\nwant:'%#v'\ngot:'%#v'\n", c.want, got)
				// }
			}

		})
	}
}
