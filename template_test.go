package late

import "testing"

func TestNewTemplate(t *testing.T) {
	tpl := NewTemplate("This is a template")

	if tpl.Body != "This is a template" {
		t.Errorf("Did not store the template body")
	}
}

func TestTemplateRender(t *testing.T) {
	tpl := NewTemplate("This is a template")
	results := tpl.Render()

	if results != "This is a template" {
		t.Errorf("Failed to render the template")
	}
}
