package json

import (
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/orznewbie/gotmpl/pkg/log"
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
