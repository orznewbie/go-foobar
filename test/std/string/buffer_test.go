package string

import (
	"bytes"
	"fmt"
	"testing"
)

func TestBytesBuffer(t *testing.T) {
	var buf bytes.Buffer
	buf.Write([]byte("Hello "))
	buf.WriteString("World!")
	fmt.Println(buf.String())
}
