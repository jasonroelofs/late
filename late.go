package late

import (
	"github.com/jasonroelofs/late/filter"
	"github.com/jasonroelofs/late/tag"
)

type tagFactoryFunc func() tag.Tag

var filters map[string]*filter.Filter
var tags map[string]tagFactoryFunc

func AddFilter(name string, filterFunc filter.FilterFunc) {
	filters[name] = filter.New(filterFunc)
}

func FindFilter(name string) *filter.Filter {
	return filters[name]
}

func AddTag(name string, tag tagFactoryFunc) {
	tags[name] = tag
}

func FindTag(name string) tag.Tag {
	tagFactory, ok := tags[name]
	if !ok {
		return nil
	}

	return tagFactory()
}

func init() {
	filters = make(map[string]*filter.Filter)
	tags = make(map[string]tagFactoryFunc)

	AddFilter("size", filter.Size)
	AddFilter("upcase", filter.Upcase)
	AddFilter("replace", filter.Replace)

	AddTag("assign", func() tag.Tag { return new(tag.Assign) })
}
