package services

import (
	"testing"

	"github.com/micro/micro/v3/service/client"
)

func TestServiceClient(t *testing.T) {
	t.Log(NewClient(client.DefaultClient))
}
