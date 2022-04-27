package time

import (
	"github.com/orznewbie/gotmpl/pkg/log"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	timer := time.NewTimer(time.Second * 3)
	for {
		log.Info("in")
		select {
		case <-timer.C:
			log.Info("Hello World!")
			continue
		}
		log.Info("xxx")
	}

	log.Info("out")
}
