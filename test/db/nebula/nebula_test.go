package nebula

import (
	nebula "github.com/vesoft-inc/nebula-go/v3"
	"testing"
)

func TestNebula(t *testing.T) {
	pool, err := nebula.NewConnectionPool()
	sess, err := pool.GetSession()
	sess.Execute()
}
