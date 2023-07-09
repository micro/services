package services

import (
	"testing"

	"micro.dev/v4/service/client"
)

func TestServiceClient(t *testing.T) {
	t.Log(NewClient(client.DefaultClient))
}
