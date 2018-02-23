package late

type Template struct {
	// The raw source of the Template we're parsing and rendering
	Body string
}

func NewTemplate(templateBody string) *Template {
	return &Template{
		Body: templateBody,
	}
}

// Render the template, returning the final output as a string
func (t *Template) Render() string {
	return t.Body
}
