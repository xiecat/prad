package output

import (
	"log"

	"github.com/tardc/prad"
)

type MultiOut struct {
	writes []Writer
}

func NewMultiOut(noColor bool, filename string) *MultiOut {
	return &MultiOut{writes: []Writer{
		NewStdout(noColor),
		NewFileOut(filename),
	}}
}

func (o *MultiOut) Write(r *prad.Result) error {
	for _, w := range o.writes {
		err := w.Write(r)
		if err != nil {
			log.Printf("write result %v failed on %v", r, w)
		}
	}

	return nil
}

func (o *MultiOut) Close() error {
	for _, w := range o.writes {
		err := w.Close()
		if err != nil {
			log.Printf("close %v failed", w)
		}
	}

	return nil
}
