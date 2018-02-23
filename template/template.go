package template

type Template struct {
	// The raw source of the Template we're parsing and rendering
	body string
}

func New(templateBody string) *Template {
	return &Template{
		body: templateBody,
	}
}

// Render the template, returning the final output as a string
func (t *Template) Render() string {
	return t.body
}
