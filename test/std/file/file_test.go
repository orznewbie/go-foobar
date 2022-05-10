package file

import (
	"io/ioutil"
	"testing"
)

func TestReadDir(t *testing.T) {
	fis, err := ioutil.ReadDir("../")
	if err != nil {
		t.Fatal(err)
	}
	for _, fi := range fis {
		t.Log(fi.Name())
	}
}
