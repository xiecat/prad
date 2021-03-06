package output

import (
	"fmt"
	"net/http"

	"github.com/logrusorgru/aurora"
	"github.com/xiecat/prad"
)

type Stdout struct {
	noColor bool
}

func NewStdout(noColor bool) *Stdout {
	return &Stdout{noColor: noColor}
}

// Write output to stdout.
func (o *Stdout) Write(r *prad.Result) error {
	var output = r.String()

	if !o.noColor {
		switch r.Code {
		case http.StatusNotFound:
			output = aurora.BrightRed(output).String()
		case http.StatusOK:
			output = aurora.BrightGreen(output).String()
		default:
			output = aurora.BrightYellow(output).String()
		}
	}

	_, err := fmt.Println(output)
	return err
}

func (o *Stdout) Close() error {
	return nil
}
