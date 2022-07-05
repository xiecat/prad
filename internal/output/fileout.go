package output

import (
	"os"

	"github.com/projectdiscovery/gologger"
	"github.com/xiecat/prad"
)

type FileOut struct {
	f *os.File
}

func NewFileOut(filename string) *FileOut {
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		gologger.Warning().Msgf("open file failed: %s", err)
		return nil
	}

	return &FileOut{f: f}
}

func (o *FileOut) Write(r *prad.Result) error {
	_, err := o.f.WriteString(r.String() + "\n")
	return err
}

func (o *FileOut) Close() error {
	return o.f.Close()
}
