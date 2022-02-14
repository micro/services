package handler

import (
	"context"

	"github.com/iverly/go-mcping/mcping"
	"github.com/micro/micro/v3/service/errors"
	pb "github.com/micro/services/minecraft/proto"
)

type Minecraft struct{}

func (m *Minecraft) Ping(ctx context.Context, req *pb.PingRequest, rsp *pb.PingResponse) error {
	if len(req.Address) == 0 {
		return errors.BadRequest("minecraft.ping", "missing address")
	}

	pinger := mcping.NewPinger()
	resp, err := pinger.Ping(req.Address, 25565)
	if err != nil {
		return err
	}

	var samples []*pb.PlayerSample
	for _, sample := range resp.Sample {
		samples = append(samples, &pb.PlayerSample{
			Uuid: sample.UUID,
			Name: sample.Name,
		})
	}

	rsp.Latency = uint32(resp.Latency)
	rsp.Players = int32(resp.PlayerCount.Online)
	rsp.MaxPlayers = int32(resp.PlayerCount.Max)
	rsp.Protocol = int32(resp.Protocol)
	rsp.Favicon = resp.Favicon
	rsp.Motd = resp.Motd
	rsp.Version = resp.Version
	rsp.Sample = samples

	return nil
}
