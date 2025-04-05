package opengraph

import "errors"

// failingWriter is an io.Writer that always returns an error.
type failingWriter struct{}

func (f *failingWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("write error")
}
