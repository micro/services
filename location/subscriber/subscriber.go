package subscriber

import (
	"context"
	"log"

	"github.com/micro/services/location/domain"
	proto "github.com/micro/services/location/proto"
)

var (
	Topic = "location"
)

type Location struct{}

func (g *Location) Handle(ctx context.Context, e *proto.Entity) error {
	log.Printf("Saving entity ID %s", e.Id)
	domain.Save(domain.ProtoToEntity(e))
	return nil
}
