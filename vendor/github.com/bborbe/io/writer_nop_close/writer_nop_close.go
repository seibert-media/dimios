package writer_nop_close

import "io"

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error {
	return nil
}

func New(r io.Writer) io.WriteCloser {
	return nopCloser{r}
}
