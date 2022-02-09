package env

const defaultTag = "env"

var defaultOptions = Options{
	Tag:      defaultTag,
	Required: false,
}

type Options struct {
	Tag      string // default "env"
	Required bool   // default false
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
