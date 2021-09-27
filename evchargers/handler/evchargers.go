package handler

import (
	"context"
	"strconv"
	"time"

	"github.com/micro/micro/v3/service/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/micro/micro/v3/service/config"
	log "github.com/micro/micro/v3/service/logger"

	evchargers "github.com/micro/services/evchargers/proto"
)

const (
	defaultDistance = int64(5000) // 5km
	//sphereIndexVersion = int32(3)
)

var (
	sphereIndexVersion = int32(3)
)

type Evchargers struct {
	conf conf
	mdb  *mongo.Client
}

type conf struct {
	MongoHost string `json:"mongo_host"`
}

func New() *Evchargers {
	val, err := config.Get("micro.evchargers")
	if err != nil {
		log.Fatalf("Failed to load config")
	}
	var conf conf
	if err := val.Scan(&conf); err != nil {
		log.Fatalf("Failed to load config")
	}
	if len(conf.MongoHost) == 0 {
		log.Fatalf("Missing mongodb host")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.MongoHost))
	if err != nil {
		log.Fatalf("Failed to connect to mongo db")
	}
	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()
	if err := client.Ping(ctx2, nil); err != nil {
		log.Fatalf("Failed to ping mongo db")
	}

	// make sure the indexes are set up
	_, err = client.Database("ocm").Collection("poi").Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{ // bounding box queries
			Keys: bson.D{{"SpatialPosition.coordinates", 1}},
			Options: &options.IndexOptions{
				SphereVersion: &sphereIndexVersion,
			},
		},
		{ // distance queries
			Keys: bson.D{{"SpatialPosition", 1}},
			Options: &options.IndexOptions{
				SphereVersion: &sphereIndexVersion,
			},
		},
		{
			Keys: bson.M{"DateCreated": -1},
		},
		{
			Keys: bson.M{"DateLastStatusUpdate": -1},
		},
		{
			Keys: bson.M{"ID": -1},
		},
	})
	if err != nil {
		log.Fatalf("Failed to craete indexes %s", err)
	}

	return &Evchargers{conf: conf, mdb: client}
}

func (e *Evchargers) Search(ctx context.Context, request *evchargers.SearchRequest, response *evchargers.SearchResponse) error {

	toInt := func(in []string) []interface{} {
		res := make([]interface{}, len(in))
		for i, v := range in {
			res[i], _ = strconv.Atoi(v)
		}
		return res
	}
	filters := bson.D{}
	if len(request.ConnectionTypes) > 0 {
		vals := bson.A{}
		vals = append(vals, toInt(request.ConnectionTypes)...)
		filters = append(filters, bson.E{"Connections.ConnectionTypeID", bson.D{{"$in", vals}}})
	}

	if request.Location != nil {
		distance := defaultDistance
		if request.Distance > 0 {
			distance = request.Distance
		}
		filters = append(filters, bson.E{"SpatialPosition", bson.M{"$nearSphere": bson.M{"$geometry": bson.M{
			"type":        "Point",
			"coordinates": []float64{float64(request.Location.Longitude), float64(request.Location.Latitude)},
		},
			"$maxDistance": distance,
		},
		}})
	}

	maxLim := int64(100)
	max := options.FindOptions{
		Limit: &maxLim,
	}
	if request.MaxResults > 0 {
		max.Limit = &request.MaxResults
	}
	crs, err := e.mdb.Database("ocm").Collection("poi").Find(ctx, filters, &max)
	if err != nil {
		log.Errorf("Error querying %s", err)
		return errors.InternalServerError("evchargers.search", "Failed to query ev chargers")
	}
	defer crs.Close(ctx)
	for crs.Next(ctx) {
		var result Poi
		if err := crs.Decode(&result); err != nil {
			log.Errorf("Error decoding result %s", err)
			return errors.InternalServerError("evchargers.search", "Failed to query ev chargers")
		}
		poi := &evchargers.Poi{
			Id:             strconv.Itoa(int(result.ID)),
			DataProviderId: strconv.Itoa(int(result.DataProviderID)),
			OperatorId:     strconv.Itoa(int(result.OperatorID)),
			UsageTypeId:    strconv.Itoa(int(result.UsageTypeID)),
			Address: &evchargers.Address{
				Location: &evchargers.Coordinates{
					Latitude:  float32(result.Address.Latitude),
					Longitude: float32(result.Address.Longitude),
				},
				Title: result.Address.Title,
			},
			Connections: marshalConnections(result.Connections),
			NumPoints:   result.NumberOfPoints,
			Cost:        result.Cost,
		}
		if true { // verbose
			poi.DataProvider = marshalDataProvider(result.DataProvider)
			poi.Operator = marshalOperator(result.OperatorInfo)
			poi.UsageType = marshalUsageType(result.UsageType)
		}
		response.Pois = append(response.Pois, poi)
	}
	return nil
}

func marshalConnections(in []Connection) []*evchargers.Connection {
	res := make([]*evchargers.Connection, len(in))
	for i, v := range in {
		res[i] = &evchargers.Connection{
			ConnectionTypeId: strconv.Itoa(int(v.TypeID)),
			ConnectionType: &evchargers.ConnectionType{
				Id:             strconv.Itoa(int(v.Type.ID)),
				Title:          v.Type.Title,
				FormalName:     v.Type.FormalName,
				IsDiscontinued: v.Type.IsDiscontinued,
				IsObsolete:     v.Type.IsObsolete,
			},
			Reference: v.Reference,
			Status:    strconv.Itoa(int(v.StatusTypeID)),
			Level:     strconv.Itoa(int(v.LevelID)),
			Amps:      float32(v.Amps),
			Voltage:   float32(v.Voltage),
			Power:     float32(v.Power),
			Current:   strconv.Itoa(int(v.CurrentTypeID)),
		}
	}
	return res
}

func marshalDataProvider(in DataProvider) *evchargers.DataProvider {
	return &evchargers.DataProvider{
		Id:                     strconv.Itoa(int(in.ID)),
		Title:                  in.Title,
		Website:                in.WebsiteURL,
		Comments:               in.Comments,
		DataProviderStatusType: marshalDataProviderStatus(in.DataProviderStatus),
		License:                in.License,
	}
}

func marshalDataProviderStatus(in DataProviderStatus) *evchargers.DataProviderStatusType {
	return &evchargers.DataProviderStatusType{
		Id:                strconv.Itoa(int(in.ID)),
		Title:             in.Title,
		IsProviderEnabled: in.IsProviderEnabled,
	}
}

func marshalOperator(in Operator) *evchargers.Operator {
	return &evchargers.Operator{
		Id:                  strconv.Itoa(int(in.ID)),
		Title:               in.Title,
		Website:             in.WebsiteURL,
		Comments:            in.Comments,
		IsPrivateIndividual: in.IsPrivateIndividual,
		ContactEmail:        in.ContactEmail,
		PhonePrimary:        in.PhonePrimary,
		PhoneSecondary:      in.PhoneSecondary,
		FaultReportEmail:    in.FaultReportEmail,
	}
}

func marshalUsageType(in UsageType) *evchargers.UsageType {
	return &evchargers.UsageType{
		Id:                   strconv.Itoa(int(in.ID)),
		Title:                in.Title,
		IsPayAtLocation:      in.IsPayAtLocation,
		IsMembershipRequired: in.IsMembershipRequired,
		IsAccessKeyRequired:  in.IsAccessKeyRequired,
	}
}

func (e *Evchargers) ReferenceData(ctx context.Context, request *evchargers.ReferenceDataRequest, response *evchargers.ReferenceDataResponse) error {
	panic("implement me")
}
