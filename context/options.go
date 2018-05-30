package context

// Self-referential functions to implement the various options
// that are available for configuring a Context

func Reader(fs FileReader) func(*Context) {
	return func(ctx *Context) {
		ctx.reader = fs
	}
}
