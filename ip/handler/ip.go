package handler

import (
	"context"
	"net"

	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	pb "github.com/micro/services/ip/proto"
	geoip2 "github.com/oschwald/geoip2-golang"
)

type Ip struct {
	ASNReader  *geoip2.Reader
	CityReader *geoip2.Reader
}

func (i *Ip) Lookup(ctx context.Context, req *pb.LookupRequest, rsp *pb.LookupResponse) error {
	if len(req.Ip) == 0 {
		return errors.BadRequest("ip.lookup", "missing ip")
	}

	// get the ip
	ip := net.ParseIP(req.Ip)

	// only if the asn reader exists
	if i.ASNReader != nil {
		asn, err := i.ASNReader.ASN(ip)
		if err != nil {
			logger.Errorf("failed to lookup asn for %v: %v", req.Ip, err)
			return errors.InternalServerError("ip.lookup", "failed to lookup ip")
		}
		// set asp
		rsp.Asn = int32(asn.AutonomousSystemNumber)
	}

	info, err := i.CityReader.City(ip)
	if err != nil {
		logger.Errorf("failed to lookup city for %v: %v", req.Ip, err)
		return errors.InternalServerError("ip.lookup", "failed to lookup ip")
	}

	// set ip
	rsp.Ip = req.Ip
	// set city
	rsp.City = info.City.Names["en"]
	// set countr
	rsp.Country = info.Country.Names["en"]
	// set continent
	rsp.Continent = info.Continent.Names["en"]
	// latitude/longitude
	rsp.Latitude = info.Location.Latitude
	rsp.Longitude = info.Location.Longitude
	// set timezone
	rsp.Timezone = info.Location.TimeZone

	return nil
}
