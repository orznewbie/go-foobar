package bloomfilter

import (
	"github.com/bits-and-blooms/bloom/v3"
	"testing"
)

func TestBasic(t *testing.T) {
	f := bloom.New(1000, 4)
	n1 := []byte("Bess")
	n2 := []byte("Jane")
	n3 := []byte("Emma")
	f.Add(n1)
	n3a := f.TestAndAdd(n3)
	n1b := f.Test(n1)
	n2b := f.Test(n2)
	n3b := f.Test(n3)
	if !n1b {
		t.Errorf("%v should be in.", n1)
	}
	if n2b {
		t.Errorf("%v should not be in.", n2)
	}
	if n3a {
		t.Errorf("%v should not be in the first time we look.", n3)
	}
	if !n3b {
		t.Errorf("%v should be in the second time we look.", n3)
	}
}
