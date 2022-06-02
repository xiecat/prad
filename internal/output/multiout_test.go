package output

import (
	"github.com/tardc/prad"
	"testing"
)

var filename = "output.txt"

func TestNewMultiOut(t *testing.T) {
	o := NewMultiOut(false, filename)
	if o == nil {
		t.Fatal("NewMultiOut failed")
	}

	err := o.Close()
	if err != nil {
		t.Fatal("MultiOut close failed")
	}
}

func TestMultiOut_Write(t *testing.T) {
	o := NewMultiOut(false, filename)
	err := o.Write(&prad.Result{
		URL:      "https://github.com/test",
		Code:     200,
		Redirect: "",
	})
	if err != nil {
		t.Fatalf("MultiOut write failed")
	}

	err = o.Close()
	if err != nil {
		t.Fatal("MultiOut close failed")
	}
}
