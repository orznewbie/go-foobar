package string

import (
	"strings"
	"testing"
)

func BenchmarkJoin(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var str []string
		for i := 0; i < 100; i++ {
			strings.Join(str, "hello")
		}
	}
}

func BenchmarkAdd(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var str string
		for i := 0; i < 100; i++ {
			str += "hello"
		}
	}
}
