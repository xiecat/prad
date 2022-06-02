package output

import (
	"github.com/tardc/prad"
)

type Writer interface {
	Write(r *prad.Result) error
	Close() error
}
