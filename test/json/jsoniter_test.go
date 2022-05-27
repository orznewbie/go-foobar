package json

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/orznewbie/gotmpl/pkg/log"
	"testing"
)

func TestJsoniter(t *testing.T) {
	byt := `[
      {
        "count": 100
      }
    ]`
	count := jsoniter.Get([]byte(byt), 0, "count").ToInt64()
	log.Info(count)
}
