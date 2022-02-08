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
				return env.AsString("INVALID")
			},
			err:  fmt.Errorf("'INVALID' not found"),
			want: nil,
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
			}

		})
	}
}
