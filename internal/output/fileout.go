package output

import (
	"github.com/tardc/prad"
	"log"
	"os"
)

type FileOut struct {
	f *os.File
}

func NewFileOut(filename string) *FileOut {
	f, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Println("open file failed:", err)
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
