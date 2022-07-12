package google_api_go_client

import (
	"context"
	"testing"

	longrunning "cloud.google.com/go/longrunning/autogen"
)

func TestAutogen(t *testing.T) {
	ctx := context.Background()
	c, err := longrunning.NewOperationsClient(ctx)
	if err != nil {
		// TODO: Handle error.
	}
	defer c.Close()
}
