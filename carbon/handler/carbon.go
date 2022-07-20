package handler

import (
	"context"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/services/carbon/domain"
	"github.com/micro/services/pkg/api"

	pb "github.com/micro/services/carbon/proto"
)

type Carbon struct {
	apiKey     string
	apiAddress string
}

func New() *Carbon {
	v, err := config.Get("carbon.api_key")
	if err != nil {
		logger.Fatalf("carbon.api_key config not found: %v", err)
	}
	apiKey := v.String("")
	if len(apiKey) == 0 {
		logger.Fatal("carbon.api_key config not found")
	}

	api.SetKey("Authorization", "Bearer "+apiKey)
	api.SetKey("Content-Type", "application/json")

	v, err = config.Get("carbon.api_address")
	if err != nil {
		logger.Fatalf("carbon.api_address config not found: %v", err)
	}
	apiAddress := v.String("")
	if len(apiKey) == 0 {
		logger.Fatal("carbon.api_address config not found")
	}

	return &Carbon{
		apiKey:     apiKey,
		apiAddress: apiAddress,
	}
}

func (c *Carbon) Offset(ctx context.Context, req *pb.OffsetRequest, rsp *pb.OffsetResponse) error {
	var resp domain.OffsetResponse

	// currently do not support options
	r := &domain.OffsetRequest{
		Number: 1,
		Units:  "KG",
	}

	if err := api.Post(c.apiAddress+"/impact/carbon", r, &resp); err != nil {
		logger.Error("Failed to purchase offsets: ", err.Error())
		return errors.InternalServerError("carbon.offset", "failed to purchase offsets")
	}

	logger.Infof("Purchased %d %s: %v\n", r.Number, r.Units, resp)

	rsp.Units = resp.Number
	rsp.Metric = resp.Units
	rsp.Tonnes = resp.Tonnes
	//rsp.Cost = resp.Amount
	//rsp.Currency = resp.Currency
	for _, p := range resp.Projects {
		rsp.Projects = append(rsp.Projects, &pb.Project{
			Name:       p.Name,
			Percentage: p.Percentage,
			Tonnes:     p.Tonnes,
		})
	}

	return nil
}
