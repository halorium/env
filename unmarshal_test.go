package env_test

import (
	"fmt"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/halorium/env"
)

type setEnv func(t *testing.T)

func TestUnmarshal(t *testing.T) {
	cases := []struct {
		name   string
		obj    interface{}
		setEnv setEnv
		opts   env.Options
		err    error
		want   interface{}
	}{
		{
			name:   "object is nil",
			obj:    nil,
			setEnv: func(t *testing.T) {},
			opts:   env.Options{},
			err:    env.ErrInvalidType,
			want:   nil,
		},
		{
			name:   "object is non-pointer",
			obj:    struct{}{},
			setEnv: func(t *testing.T) {},
			opts:   env.Options{},
			err:    env.ErrInvalidType,
			want:   nil,
		},
		{
			name:   "object is non-struct",
			obj:    ptr(""),
			setEnv: func(t *testing.T) {},
			opts:   env.Options{},
			err:    env.ErrInvalidType,
			want:   nil,
		},
		{
			name: "valid string field",
			obj: &struct {
				String string `env:"STRING"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("STRING", "string_val")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				String string `env:"STRING"`
			}{String: "string_val"},
		},
		{
			name: "valid *string field",
			obj: &struct {
				String *string `env:"STRING"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("STRING", "string_val")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				String *string `env:"STRING"`
			}{String: ptrStr("string_val")},
		},
		{
			name: "valid bool field",
			obj: &struct {
				Bool bool `env:"BOOL"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("BOOL", "true")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				Bool bool `env:"BOOL"`
			}{Bool: true},
		},
		{
			name: "valid *bool field",
			obj: &struct {
				Bool *bool `env:"BOOL"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("BOOL", "true")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				Bool *bool `env:"BOOL"`
			}{Bool: ptrBool(true)},
		},
		{
			name: "valid int field",
			obj: &struct {
				Int int `env:"INT"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("INT", "1")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				Int int `env:"INT"`
			}{Int: 1},
		},
		{
			name: "valid *int field",
			obj: &struct {
				Int *int `env:"INT"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("INT", "1")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				Int *int `env:"INT"`
			}{Int: ptrInt(1)},
		},
		{
			name: "valid float32 field",
			obj: &struct {
				Float float32 `env:"FLOAT"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("FLOAT", "2.3")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				Float float32 `env:"FLOAT"`
			}{Float: 2.3},
		},
		{
			name: "valid *float32 field",
			obj: &struct {
				Float *float32 `env:"FLOAT"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("FLOAT", "2.3")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				Float *float32 `env:"FLOAT"`
			}{Float: ptrFloat32(2.3)},
		},
		{
			name: "valid float64 field",
			obj: &struct {
				Float float64 `env:"FLOAT"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("FLOAT", "2.3")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				Float float64 `env:"FLOAT"`
			}{Float: 2.3},
		},
		{
			name: "valid *float64 field",
			obj: &struct {
				Float *float64 `env:"FLOAT"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("FLOAT", "2.3")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				Float *float64 `env:"FLOAT"`
			}{Float: ptrFloat64(2.3)},
		},
		{
			name: "valid duration field",
			obj: &struct {
				Duration time.Duration `env:"DURATION"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("DURATION", "5s")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				Duration time.Duration `env:"DURATION"`
			}{Duration: 5 * time.Second},
		},
		{
			name: "valid *duration field",
			obj: &struct {
				Duration *time.Duration `env:"DURATION"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("DURATION", "5s")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				Duration *time.Duration `env:"DURATION"`
			}{Duration: ptrDuration(5 * time.Second)},
		},
		{
			name: "valid []string field",
			obj: &struct {
				List []string `env:"LIST"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("LIST", "one,two")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				List []string `env:"LIST"`
			}{List: []string{"one", "two"}},
		},
		{
			name: "valid *[]string field",
			obj: &struct {
				List *[]string `env:"LIST"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("LIST", "one,two")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				List *[]string `env:"LIST"`
			}{List: &[]string{"one", "two"}},
		},
		{
			name: "valid []int field",
			obj: &struct {
				List []int `env:"LIST"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("LIST", "1,2")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				List []int `env:"LIST"`
			}{List: []int{1, 2}},
		},
		{
			name: "valid *[]int field",
			obj: &struct {
				List *[]int `env:"LIST"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("LIST", "1,2")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				List *[]int `env:"LIST"`
			}{List: &[]int{1, 2}},
		},
		{
			name: "valid []int field empty env",
			obj: &struct {
				List []int `env:"LIST"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("LIST", "")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				List []int `env:"LIST"`
			}{List: []int{}},
		},
		{
			name: "valid []byte field",
			obj: &struct {
				Byte []byte `env:"BYTE"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("BYTE", "some data")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				Byte []byte `env:"BYTE"`
			}{Byte: []byte("some data")},
		},
		{
			name: "valid *[]byte field",
			obj: &struct {
				Byte *[]byte `env:"BYTE"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("BYTE", "some data")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				Byte *[]byte `env:"BYTE"`
			}{Byte: ptrByte([]byte("some data"))},
		},
		{
			name: "valid map[string]int field",
			obj: &struct {
				Map map[string]int `env:"MAP"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("MAP", "one:1,two:2")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				Map map[string]int `env:"MAP"`
			}{Map: map[string]int{"one": 1, "two": 2}},
		},
		{
			name: "valid *map[string]int field",
			obj: &struct {
				Map *map[string]int `env:"MAP"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("MAP", "one:1,two:2")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				Map *map[string]int `env:"MAP"`
			}{Map: &map[string]int{"one": 1, "two": 2}},
		},
		{
			name: "valid uint32 field",
			obj: &struct {
				Int uint32 `env:"INT"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("INT", "1")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				Int uint32 `env:"INT"`
			}{Int: 1},
		},
		{
			name: "valid *uint32 field",
			obj: &struct {
				Int *uint32 `env:"INT"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("INT", "1")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				Int *uint32 `env:"INT"`
			}{Int: ptrUint32(1)},
		},
		{
			name: "valid string field required option error",
			obj: &struct {
				String string `env:"STRING"`
			}{},
			setEnv: func(t *testing.T) {},
			opts:   env.Options{Required: true},
			err:    RequiredErr("STRING"),
			want:   nil,
		},
		{
			name: "valid string field required option no error",
			obj: &struct {
				String string `env:"STRING"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("STRING", "string_val")
			},
			opts: env.Options{Required: true},
			err:  nil,
			want: &struct {
				String string `env:"STRING"`
			}{String: "string_val"},
		},
		{
			name: "valid string field env:'-'",
			obj: &struct {
				String string `env:"-"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("STRING", "string_val")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				String string `env:"-"`
			}{},
		},
		{
			name: "valid string field env:'-' option required",
			obj: &struct {
				String string `env:"-"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("STRING", "string_val")
			},
			opts: env.Options{Required: true},
			err:  nil,
			want: &struct {
				String string `env:"-"`
			}{},
		},
		{
			name: "valid string field env:''",
			obj: &struct {
				String string `env:""`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("STRING", "string_val")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				String string `env:""`
			}{String: ""},
		},
		{
			name: "valid string field env:'' option required",
			obj: &struct {
				String string `env:""`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("STRING", "string_val")
			},
			opts: env.Options{Required: true},
			err:  nil,
			want: &struct {
				String string `env:""`
			}{String: ""},
		},
		{
			name: "valid string field on nested structs",
			obj:  &ParentStruct{},
			setEnv: func(t *testing.T) {
				t.Setenv("NESTED", "nested_val")
				t.Setenv("EMBEDDED", "embedded_val")
			},
			opts: env.Options{},
			err:  nil,
			want: &ParentStruct{
				EmbeddedStruct: EmbeddedStruct{String: "embedded_val"},
				Nested:         NestedStruct{String: "nested_val"},
				NestedPtr:      &NestedStruct{String: "nested_val"},
			},
		},
		{
			name: "valid string field on nested structs instantiated",
			obj: &ParentStruct{
				Nested:    NestedStruct{},
				NestedPtr: &NestedStruct{},
			},
			setEnv: func(t *testing.T) {
				t.Setenv("NESTED", "nested_val")
				t.Setenv("EMBEDDED", "embedded_val")
			},
			opts: env.Options{},
			err:  nil,
			want: &ParentStruct{
				EmbeddedStruct: EmbeddedStruct{String: "embedded_val"},
				Nested:         NestedStruct{String: "nested_val"},
				NestedPtr:      &NestedStruct{String: "nested_val"},
			},
		},
		{
			name: "valid custom field",
			obj: &struct {
				IP CustomIP `env:"IP"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("IP", "127.0.0.1")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				IP CustomIP `env:"IP"`
			}{IP: CustomIP("127.0.0.1")},
		},
		{
			name: "valid custom struct field with Unmarshaler",
			obj: &struct {
				URL CustomURL `env:"URL"`
			}{},
			setEnv: func(t *testing.T) {
				t.Setenv("URL", "http://github.com/halorium/env")
			},
			opts: env.Options{},
			err:  nil,
			want: &struct {
				URL CustomURL `env:"URL"`
			}{URL: getCustomURL()},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.setEnv(t)
			err := env.Unmarshal(c.obj, c.opts)
			if err != nil && c.err != nil {
				if err.Error() != c.err.Error() {
					t.Errorf("\nwant:'%#v'\ngot:'%#v'\n", c.err.Error(), err.Error())
				}
			} else if err != c.err {
				t.Errorf("\nwant:'%#v'\ngot:'%#v'\n", c.err, err)
			}

			if err == nil {
				if diff := cmp.Diff(c.want, c.obj); diff != "" {
					t.Errorf("(-want +got):\n%s", diff)
				}

				// if !reflect.DeepEqual(c.want, c.obj) {
				// 	t.Errorf("\nwant:'%#v'\ngot:'%#v'\n", c.want, c.obj)
				// }

			}
		})
	}
}

func ptr(obj interface{}) interface{} {
	return &obj
}

func ptrStr(v string) *string {
	return &v
}

func ptrBool(v bool) *bool {
	return &v
}

func ptrInt(v int) *int {
	return &v
}

func ptrFloat32(v float32) *float32 {
	return &v
}

func ptrFloat64(v float64) *float64 {
	return &v
}

func ptrDuration(v time.Duration) *time.Duration {
	return &v
}

func ptrByte(v []byte) *[]byte {
	return &v
}

func ptrUint32(v uint32) *uint32 {
	return &v
}

type ParentStruct struct {
	EmbeddedStruct
	Nested    NestedStruct
	NestedPtr *NestedStruct
}

type NestedStruct struct {
	String string `env:"NESTED"`
}

type EmbeddedStruct struct {
	String string `env:"EMBEDDED"`
}

type CustomIP net.IP

type CustomURL url.URL

func getCustomURL() CustomURL {
	u, _ := url.Parse("http://github.com/halorium/env")
	return CustomURL(*u)
}

func (c *CustomURL) UnmarshalENV(v string) error {
	u, e := url.Parse(v)
	if e != nil {
		return e
	}
	*c = CustomURL(*u)
	return nil
}

func RequiredErr(v string) error {
	return fmt.Errorf("'%s' is required", v)
}
