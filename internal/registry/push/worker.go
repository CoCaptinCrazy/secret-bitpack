package push

import "io"

type Worker struct {
	Source      io.Reader
	Destination io.Writer
}

func (w Worker) Run() (int64, error) {
	return io.Copy(w.Destination, w.Source)
}
