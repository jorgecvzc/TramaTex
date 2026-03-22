package sharedomain

// HTTPStatuser is implemented by domain errors that carry their own HTTP status code.
// The global error middleware checks for this interface to determine the response status.
type HTTPStatuser interface {
	error
	HTTPStatus() int
}
