package output

import (
	"github.com/xiecat/prad"
)

type Writer interface {
	Write(r *prad.Result) error
	Close() error
}
