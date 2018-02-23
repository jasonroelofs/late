package template

import "testing"

func TestNew(t *testing.T) {
	tpl := New("This is a template")

	if tpl.body != "This is a template" {
		t.Errorf("Did not store the template body")
	}
}

func TestRender(t *testing.T) {
	tpl := New("This is a template")
	results := tpl.Render()

	if results != "This is a template" {
		t.Errorf("Failed to render the template")
	}
}
