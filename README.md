# env

```Go
import "github.com/halorium/env"
```

## Usage

Set some environment variables:

```Bash
export TEST_STRING="some value"
export TEST_BOOL=true
export TEST_INT=1
export TEST_FLOAT=1.2
export TEST_DURATION="5s"
export TEST_LIST="one,two,three"
export TEST_MAP="one:1,two:2,three:3"
```

Write some code:

```Go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/halorium/env"
)

type Config struct {
	StringField   string         `env:"TEST_STRING"`
	BoolField     bool           `env:"TEST_BOOL"`
	IntField      int            `env:"TEST_INT"`
	FloatField    float64        `env:"TEST_FLOAT"`
	DurationField time.Duration  `env:"TEST_DURATION"`
	ListField     []string       `env:"TEST_LIST"`
	MapField      map[string]int `env:"TEST_MAP"`
}

func (c Config) String() string {
	val := fmt.Sprintf(
		"StringField:'%s'\nBoolField:'%v'\nIntField:'%d'\nFloatField:'%f'\nDurationField:'%s'\n",
		c.StringField,
		c.BoolField,
		c.IntField,
		c.FloatField,
		c.DurationField)

	val += "ListField:\n  " + strings.Join(c.ListField, "  \n") + "\n"

	val += "MapField:\n"
	for k, v := range c.MapField {
		val += fmt.Sprintf("  %s: %d\n", k, v)
	}

	return val
}

func main() {
	var cfg Config
	err := env.Unmarshal(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)
}
```

Results:

```Bash
 todo
```

## Struct Tags

The default tag is '`env`' however this can be changed in the options.

Example:
```Go
type Config struct {
	StringField   string         `custom:"TEST_STRING"`
	BoolField     bool           `custom:"TEST_BOOL"`
	IntField      int            `custom:"TEST_INT"`
	FloatField    float64        `custom:"TEST_FLOAT"`
	DurationField time.Duration  `custom:"TEST_DURATION"`
	ListField     []string       `custom:"TEST_LIST"`
	MapField      map[string]int `custom:"TEST_MAP"`
}

var cfg Config
options := env.Options{Tag: "custom"}
err := env.Unmarshal(&cfg, options)
if err != nil {
	log.Fatal(err)
}
fmt.Println(cfg)

```

## Validation
env doesn't assume any validation, if an environment variable is not found then it is skipped.

The '`Required`' option can be set to have an error returned if any environment variables are not set.
Example:
```go
var cfg Config
options := env.Options{Required: true}
err := env.Unmarshal(&cfg, options)
if err != nil {
	log.Fatal(err)
}
fmt.Println(cfg)
```

## Ignored Fields
env will ignore the field if the tag is set to either an empty string '' or a hyphen '-'.
Example:
```go
type Config struct {
	StringField   string `env:"TEST_STRING"`
	BoolField     bool   `env:""`  // ignored
	IntField      int    `env:"-"` // ignored
}
```

## Supported Field Types

* string
* int8, int16, int32, int64
* bool
* float32, float64
* slices of any supported type
* maps (keys and values of any supported type)
* [time.Duration](https://golang.org/pkg/time/#Duration)


`Note`: Embedded structs using these fields are also supported.

## Custom Unmarshaler

Any field whose type (or pointer-to-type) implements `env.UnmarshalENV` can
control its own deserialization:

```Bash
export IP=127.0.0.1
```

```Go
type CustomIP net.IP

func (c *CustomIP) UnmarshalENV(v string) error {
    *c = CustomIP(net.ParseIP(v))
    return nil
}

type Config struct {
    IP CustomIP `env:"IP"`
}
```

## Other Ways to Use '`env`'

env implements helper functions if you simply want to get environment variables as a specified type.

```go
s, err := env.AsString("TEST_STRING") // returns (string, error)

b, err := env.AsBool("TEST_BOOL") // returns (bool, error)

i, err := env.AsInt("TEST_INT", 64) // returns (int64, error)

i, err := env.AsUint("TEST_UINT", 64) // returns (uint64, error)

f, err := env.AsFloat("TEST_FLOAT", 64) // returns (float64, error)

d, err := env.AsDuration("TEST_DURATION") // returns (time.Duration, error)
```