[![Build Status](https://app.travis-ci.com/halorium/env.svg?branch=main)](https://app.travis-ci.com/halorium/env)

# env

## A Zero Dependency - No Nonsense Environment Variable Parser / Unmarshaler

```Go
import "github.com/halorium/env"
```

## Usage

### There are two ways to use `env`

* Unmarshal into a struct using struct tags
* Get values as specified type for use how you want (no struct needed)

### Example
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

func main() {
	var cfg Config
	err := env.Unmarshal(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(cfg)
}
```

## Other Ways to Use '`env`'

Access environment variables as specified types directly.

```go
s, err := env.AsString("TEST_STRING") // returns (string, error)

b, err := env.AsBool("TEST_BOOL") // returns (bool, error)

i, err := env.AsInt("TEST_INT", 64) // returns (int64, error)

i, err := env.AsUint("TEST_UINT", 64) // returns (uint64, error)

f, err := env.AsFloat("TEST_FLOAT", 64) // returns (float64, error)

d, err := env.AsDuration("TEST_DURATION") // returns (time.Duration, error)
```
## Options (Struct Tags, Validation, etc.)

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
`env` doesn't assume any validation, if an environment variable is not found then it is skipped.

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
`env` will ignore the field if the tag is set to either an empty string '' or a hyphen '-'.
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
* bool
* int, int8, int16, int32, int64
* uint, uint8, uint16, uint32, uint64
* float32, float64
* slices of any supported type
* maps (keys and values of any supported type)
* [time.Duration](https://golang.org/pkg/time/#Duration)
* any field that implements the Unmarshaler interface (UnmarshalENV)


`Note`: Embedded structs using these fields are also supported.

## Custom Unmarshaler

Any field whose type (or pointer-to-type) implements `env.UnmarshalENV` can
control its own deserialization:

```Bash
export URL="http://github.com/halorium/env"
```

```Go
type CustomURL url.URL

func (c *CustomURL) UnmarshalENV(v string) error {
	url, err := url.Parse(v)
	if err != nil {
		return err
	}
	*c = CustomURL(*url)
	return nil
}

type Config struct {
    URL CustomURL `env:"URL"`
}
```
