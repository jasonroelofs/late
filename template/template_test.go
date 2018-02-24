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

//func TestRenderWithParameters(t *testing.T) {
//	tpl := New("This is a {{ thing }}")
//	results := tpl.Render(Params{"thing": "template"})
//
//	if results != "This is a template" {
//		t.Errorf("Failed to render the template. Expected '%s' got '%s'", "This is a template", results)
//	}
//}
