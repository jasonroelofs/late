package context

import (
	"github.com/jasonroelofs/late/object"
)

type Scope struct {
	Parent  *Scope
	assigns map[string]object.Object
	shadow  bool
}

func NewScope(parent *Scope) *Scope {
	return &Scope{
		Parent:  parent,
		assigns: make(map[string]object.Object),
	}
}

func NewShadowScope(parent *Scope) *Scope {
	newScope := NewScope(parent)
	newScope.shadow = true
	return newScope
}

func (s *Scope) Set(name string, value object.Object) {
	if s.shadow {
		// It's not possible to have a shadow root scope, so we don't check
		// for that case here. I expect to be proven wrong at some point.
		s.Parent.Set(name, value)
		return
	}

	s.assigns[name] = value
}

func (s *Scope) ShadowSet(name string, value object.Object) {
	if !s.shadow {
		return
	}

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
