package context

type FileReader interface {
	// Given a path, return the content of the file
	// at that path. Use custom readers to define how
	// to find other templates and partials (particularly with the
	// `include` tag).
	Read(string) string
}

// Provide a default null file system implementation
// that always returns an error
type NullReader struct{}

func (n *NullReader) Read(path string) string {
	return "ERROR: Reader not implemented. Cannot read content at " + path
}
