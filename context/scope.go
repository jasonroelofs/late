package context

import (
	"github.com/jasonroelofs/late/object"
)

type Scope struct {
	Parent  *Scope
	assigns map[string]object.Object
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		Parent:  parent,
		assigns: make(map[string]object.Object),
	}
}

func (s *Scope) Set(name string, value object.Object) {
	s.assigns[name] = value
}

func (s *Scope) Get(name string) object.Object {
	obj, ok := s.assigns[name]

	if !ok {
		if s.Parent == nil {
			return object.NULL
		} else {
			return s.Parent.Get(name)
		}
	}

	return obj
}

func (s *Scope) Promote(name string) {
	// Doesn't make sense, but lets not error out.
	// The variable is already promoted as far as it can go.
	if s.Parent == nil {
		return
	}

	// TODO: If this variable doesn't exist in current scope?
	s.Parent.Set(name, s.Get(name))
}
