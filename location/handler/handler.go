package handler

import (
	"context"
	"log"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/location/domain"
	loc "github.com/micro/services/location/proto"
	"github.com/micro/services/location/subscriber"
	pauth "github.com/micro/services/pkg/auth"
	adminpb "github.com/micro/services/pkg/service/proto"
)

type Location struct{}

func (l *Location) Read(ctx context.Context, req *loc.ReadRequest, rsp *loc.ReadResponse) error {
	log.Print("Received Location.Read request")

	id := req.Id

	if len(id) == 0 {
		return errors.BadRequest("location.read", "Require Id")
	}

	entity, err := domain.Read(ctx, id)
	if err != nil {
		return err
	}

	rsp.Entity = entity.ToProto()

	return nil
}

func (l *Location) Save(ctx context.Context, req *loc.SaveRequest, rsp *loc.SaveResponse) error {
	log.Print("Received Location.Save request")

	entity := req.GetEntity()

	if entity.GetLocation() == nil {
		return errors.BadRequest("location.save", "Require location")
	}

	// immediate save
	domain.Save(ctx, domain.ProtoToEntity(entity))

	// publish the event so other copies of location service can save it
	p := service.NewEvent(subscriber.Topic)

	if err := p.Publish(ctx, entity); err != nil {
		return errors.InternalServerError("location.save", err.Error())
	}

	return nil
}

func (l *Location) Search(ctx context.Context, req *loc.SearchRequest, rsp *loc.SearchResponse) error {
	log.Print("Received Location.Search request")

	entity := &domain.Entity{
		Latitude:  req.Center.Latitude,
		Longitude: req.Center.Longitude,
	}

	entities := domain.Search(ctx, req.Type, entity, req.Radius, int(req.Limit))

	for _, e := range entities {
		rsp.Entities = append(rsp.Entities, e.ToProto())
	}

	return nil
}

func (l *Location) DeleteData(ctx context.Context, request *adminpb.DeleteDataRequest, response *adminpb.DeleteDataResponse) error {
	method := "admin.DeleteData"
	_, err := pauth.VerifyMicroAdmin(ctx, method)
	if err != nil {
		return err
	}

	if len(request.TenantId) < 10 { // deliberate length check so we don't delete all the things
		return errors.BadRequest(method, "Missing tenant ID")
	}

	if err := domain.DeleteIndex(request.TenantId); err != nil {
		return err
	}
	logger.Infof("Deleted index for %s", request.TenantId)
	return nil
}

func (l *Location) Usage(ctx context.Context, request *adminpb.UsageRequest, response *adminpb.UsageResponse) error {
	return nil
}
