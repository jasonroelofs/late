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

//func TestRenderLiquidWithLiterals(t *testing.T) {
//	tests := []struct {
//		liquidInput    string
//		expectedOutput string
//	}{
//		{"{{ 3 }}", "3"},
//		{"{{ 1 + 2 }}", "3"},
//		{"{{ 1 - 2 }}", "-2"},
//		{"{{ 1 / 2 }}", "0.5"},
//		{"{{ \"Hi\" }}", "Hi"},
//		{"{{ 'Hi' + ' ' + 'Bye' }}", "Hi Bye"},
//	}
//
//	for _, test := range tests {
//		tpl := New(test.liquidInput)
//		results := tpl.Render()
//
//		if results != test.expectedOutput {
//			t.Errorf("Failed to render the template. Expected '%s' got '%s'", test.expectedOutput, results)
//		}
//	}
//}
