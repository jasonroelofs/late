package statement

/**
 * A generic interface for passing around Statement nodes between the
 * context, evaluator, and tags.
 *
 * Split into its own package to prevent cyclical import errors.
 */
type Statement interface {
	String() string
}
