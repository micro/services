package handler

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/micro/micro/v3/service/errors"
	"github.com/robfig/cron/v3"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/micro/micro/v3/service/config"
	log "github.com/micro/micro/v3/service/logger"

	evchargers "github.com/micro/services/evchargers/proto"
)

const (
	defaultDistance = int64(5000) // 5km

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
	CaCrt     string `json:"ca_crt"`
	OCMKey    string `json:"ocm_key"`
}

func New() *Evchargers {
	val, err := config.Get("evchargers")
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
	if len(conf.CaCrt) > 0 {
		// write the cert to file
		if err := ioutil.WriteFile(os.TempDir()+"/mongo.crt", []byte(conf.CaCrt), 0644); err != nil {
			log.Fatalf("Failed to write crt file for mongodb %s", err)
		}
	}
	opts := []*options.ClientOptions{options.Client().ApplyURI(conf.MongoHost)}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, opts...)
	if err != nil {
		log.Fatalf("Failed to connect to mongo db %s", err)
	}
	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()
	if err := client.Ping(ctx2, nil); err != nil {
		log.Fatalf("Failed to ping mongo db %s", err)
	}

	// make sure the indexes are set up
	_, err = client.Database("ocm").Collection("poi").Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{ // bounding box queries
			Keys: bson.D{{"SpatialPosition.coordinates", "2dsphere"}},
			Options: &options.IndexOptions{
				SphereVersion: &sphereIndexVersion,
			},
		},
		{ // distance queries
			Keys: bson.D{{"SpatialPosition", "2dsphere"}},
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

	ev := &Evchargers{conf: conf, mdb: client}
	if len(conf.OCMKey) > 0 {
		c := cron.New()
		// 4am every Sunday for refresh
		c.AddFunc("0 4 * * 0", ev.refreshDataFromSource)
		c.Start()
	}

	return ev
}

func (e *Evchargers) Search(ctx context.Context, request *evchargers.SearchRequest, response *evchargers.SearchResponse) error {

	addFilter := func(filters bson.D, key, op string, in []string) bson.D {
		vals := bson.A{}
		for _, v := range in {
			if v == "" {
				continue
			}
			r, _ := strconv.Atoi(v)
			vals = append(vals, r)
		}
		if len(vals) == 0 {
			return filters
		}

		filters = append(filters, bson.E{key, bson.D{{op, vals}}})
		return filters
	}
	filters := bson.D{}

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
	} else if request.Box != nil && request.Box.BottomLeft != nil {
		filters = append(filters, bson.E{"SpatialPosition.coordinates", bson.M{"$geoWithin": bson.M{"$box": bson.A{
			[]float32{request.Box.BottomLeft.Longitude, request.Box.BottomLeft.Latitude},
			[]float32{request.Box.TopRight.Longitude, request.Box.TopRight.Latitude},
		}}}})
	}

	if len(request.CountryId) > 0 {
		i, _ := strconv.Atoi(request.CountryId)
		filters = append(filters, bson.E{"AddressInfo.CountryID", i})
	}

	if len(request.ConnectionTypes) > 0 {
		filters = addFilter(filters, "Connections.ConnectionTypeID", "$in", request.ConnectionTypes)
	}

	if len(request.Levels) > 0 {
		filters = addFilter(filters, "Connections.LevelID", "$in", request.Levels)
	}

	if request.MinPower > 0 {
		filters = append(filters, bson.E{"Connections.PowerKW", bson.D{{"$gte", request.MinPower}}})
	}

	if len(request.Operators) > 0 {
		filters = addFilter(filters, "OperatorID", "$in", request.Operators)
	}

	if len(request.UsageTypes) > 0 {
		filters = addFilter(filters, "UsageTypeID", "$in", request.UsageTypes)
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
				Title:           result.Address.Title,
				AddressLine_1:   result.Address.AddressLine1,
				AddressLine_2:   result.Address.AddressLine2,
				Town:            result.Address.Town,
				StateOrProvince: result.Address.StateOrProvince,
				AccessComments:  result.Address.AccessComments,
				Postcode:        result.Address.Postcode,
				CountryId:       strconv.Itoa(int(result.Address.CountryID)),
				LatLng:          fmt.Sprintf("%v, %v", result.Address.Latitude, result.Address.Longitude),
			},
			Connections: marshalConnections(result.Connections),
			NumPoints:   int64(result.NumberOfPoints),
			Cost:        result.Cost,
		}
		if true { // verbose
			poi.Operator = marshalOperator(result.OperatorInfo)
			poi.UsageType = marshalUsageType(result.UsageType)
			poi.Address.Country = marshalCountry(result.Address.Country)
		}
		response.Pois = append(response.Pois, poi)
	}
	return nil
}

func marshalCountry(in Country) *evchargers.Country {
	return &evchargers.Country{
		Id:            strconv.Itoa(int(in.ID)),
		Title:         in.Title,
		IsoCode:       in.ISOCode,
		ContinentCode: in.ContinentCode,
	}
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
			LevelId:   strconv.Itoa(int(v.LevelID)),
			Level: &evchargers.ChargerType{
				Id:                  strconv.Itoa(int(v.Level.ID)),
				Title:               v.Level.Title,
				Comments:            v.Level.Comments,
				IsFastChargeCapable: v.Level.IsFastChargeCapable,
			},
			Amps:    float32(v.Amps),
			Voltage: float32(v.Voltage),
			Power:   float32(v.Power),
			Current: strconv.Itoa(int(v.CurrentTypeID)),
			Status:  marshalStatusType(v.StatusType),
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

func marshalCheckinStatusType(in CheckinStatusType) *evchargers.CheckinStatusType {
	return &evchargers.CheckinStatusType{
		Id:          strconv.Itoa(int(in.ID)),
		Title:       in.Title,
		IsPositive:  in.IsPositive,
		IsAutomated: in.IsAutomatedCheckin,
	}
}

func marshalUserCommentType(in UserCommentType) *evchargers.UserCommentType {
	return &evchargers.UserCommentType{
		Id:    strconv.Itoa(int(in.ID)),
		Title: in.Title,
	}
}

func marshalStatusType(in StatusType) *evchargers.StatusType {
	return &evchargers.StatusType{
		Id:            strconv.Itoa(int(in.ID)),
		Title:         in.Title,
		IsOperational: in.IsOperational,
	}
}

func marshalCurrentType(in CurrentType) *evchargers.CurrentType {
	return &evchargers.CurrentType{
		Id:          strconv.Itoa(int(in.ID)),
		Title:       in.Title,
		Description: in.Description,
	}
}

func marshalConnectionType(in ConnectionType) *evchargers.ConnectionType {
	return &evchargers.ConnectionType{
		Id:             strconv.Itoa(int(in.ID)),
		Title:          in.Title,
		FormalName:     in.FormalName,
		IsDiscontinued: in.IsDiscontinued,
		IsObsolete:     in.IsObsolete,
	}
}

func marshalChargerType(in ChargerType) *evchargers.ChargerType {
	return &evchargers.ChargerType{
		Id:                  strconv.Itoa(int(in.ID)),
		Title:               in.Title,
		Comments:            in.Comments,
		IsFastChargeCapable: in.IsFastChargeCapable,
	}
}

func marshalSubmissionStatusType(in SubmissionStatusType) *evchargers.SubmissionStatusType {
	return &evchargers.SubmissionStatusType{
		Id:     strconv.Itoa(int(in.ID)),
		Title:  in.Title,
		IsLive: in.IsLive,
	}
}

func (e *Evchargers) ReferenceData(ctx context.Context, request *evchargers.ReferenceDataRequest, response *evchargers.ReferenceDataResponse) error {

	res := e.mdb.Database("ocm").Collection("reference").FindOne(ctx, bson.D{})
	if res.Err() != nil {
		log.Errorf("Error retrieving ref data %s", res.Err())
		return errors.InternalServerError("evchargers.referencedata", "Error retrieving reference data")
	}
	var r ReferenceData
	if err := res.Decode(&r); err != nil {
		log.Errorf("Error decoding ref data %s", err)
		return errors.InternalServerError("evchargers.referencedata", "Error retrieving reference data")
	}
	dps := make([]*evchargers.DataProvider, len(r.DataProviders))
	for i, dp := range r.DataProviders {
		dps[i] = marshalDataProvider(dp)
	}
	response.DataProviders = dps
	cs := make([]*evchargers.Country, len(r.Countries))
	for i, c := range r.Countries {
		cs[i] = marshalCountry(c)
	}
	response.Countries = cs

	cst := make([]*evchargers.CheckinStatusType, len(r.CheckinStatusTypes))
	for i, v := range r.CheckinStatusTypes {
		cst[i] = marshalCheckinStatusType(v)
	}
	response.CheckinStatusTypes = cst

	uct := make([]*evchargers.UserCommentType, len(r.UserCommentTypes))
	for i, v := range r.UserCommentTypes {
		uct[i] = marshalUserCommentType(v)
	}
	response.UserCommentTypes = uct

	st := make([]*evchargers.StatusType, len(r.StatusTypes))
	for i, v := range r.StatusTypes {
		st[i] = marshalStatusType(v)
	}
	response.StatusTypes = st

	ut := make([]*evchargers.UsageType, len(r.UsageTypes))
	for i, v := range r.UsageTypes {
		ut[i] = marshalUsageType(v)
	}
	response.UsageTypes = ut

	ct := make([]*evchargers.CurrentType, len(r.CurrentTypes))
	for i, v := range r.CurrentTypes {
		ct[i] = marshalCurrentType(v)
	}
	response.CurrentTypes = ct

	connt := make([]*evchargers.ConnectionType, len(r.ConnectionTypes))
	for i, v := range r.ConnectionTypes {
		connt[i] = marshalConnectionType(v)
	}
	response.ConnectionTypes = connt
	chrgt := make([]*evchargers.ChargerType, len(r.ChargerTypes))
	for i, v := range r.ChargerTypes {
		chrgt[i] = marshalChargerType(v)
	}
	response.ChargerTypes = chrgt

	ops := make([]*evchargers.Operator, len(r.Operators))
	for i, v := range r.Operators {
		ops[i] = marshalOperator(v)
	}
	response.Operators = ops

	sst := make([]*evchargers.SubmissionStatusType, len(r.SubmissionStatusTypes))
	for i, v := range r.SubmissionStatusTypes {
		sst[i] = marshalSubmissionStatusType(v)
	}

	response.SubmissionStatusTypes = sst
	return nil
}
