package domain

import (
	"context"
	"sync"

	geo "github.com/hailocab/go-geoindex"
	"github.com/micro/micro/v3/service/errors"
	common "github.com/micro/services/location/proto"
	"github.com/micro/services/pkg/tenant"
)

var (
	mtx          sync.RWMutex
	defaultIndex = geo.NewPointsIndex(geo.Km(100))

	// index per tenant
	indexes = map[string]*geo.PointsIndex{}
)

type Entity struct {
	ID        string
	Type      string
	Latitude  float64
	Longitude float64
	Timestamp int64
}

func (e *Entity) Id() string {
	return e.ID
}

func (e *Entity) Lat() float64 {
	return e.Latitude
}

func (e *Entity) Lon() float64 {
	return e.Longitude
}

func (e *Entity) ToProto() *common.Entity {
	return &common.Entity{
		Id:   e.ID,
		Type: e.Type,
		Location: &common.Point{
			Latitude:  e.Latitude,
			Longitude: e.Longitude,
			Timestamp: e.Timestamp,
		},
	}
}

func getIndex(ctx context.Context) *geo.PointsIndex {
	tenantId, ok := tenant.FromContext(ctx)
	if !ok {
		//return default index
		return defaultIndex
	}

	// get the index
	index, ok := indexes[tenantId]
	if !ok {
		index = geo.NewPointsIndex(geo.Meters(100))
		indexes[tenantId] = index
	}

	return index
}

func ProtoToEntity(e *common.Entity) *Entity {
	return &Entity{
		ID:        e.Id,
		Type:      e.Type,
		Latitude:  e.Location.Latitude,
		Longitude: e.Location.Longitude,
		Timestamp: e.Location.Timestamp,
	}
}

func Read(ctx context.Context, id string) (*Entity, error) {
	mtx.RLock()
	defer mtx.RUnlock()

	// get the index
	index := getIndex(ctx)

	p := index.Get(id)
	if p == nil {
		return nil, errors.NotFound("location.read", "Not found")
	}

	entity, ok := p.(*Entity)
	if !ok {
		return nil, errors.InternalServerError("location.read", "Error reading entity")
	}

	return entity, nil
}

func Save(ctx context.Context, e *Entity) {
	mtx.Lock()
	// get the index
	index := getIndex(ctx)
	index.Add(e)
	mtx.Unlock()
}

func Search(ctx context.Context, typ string, entity *Entity, radius float64, numEntities int) []*Entity {
	mtx.RLock()
	defer mtx.RUnlock()

	// get the index
	index := getIndex(ctx)

	points := index.KNearest(entity, numEntities, geo.Meters(radius), func(p geo.Point) bool {
		e, ok := p.(*Entity)
		if !ok || e.Type != typ {
			return false
		}
		return true
	})

	var entities []*Entity

	for _, point := range points {
		e, ok := point.(*Entity)
		if !ok {
			continue
		}
		entities = append(entities, e)
	}

	return entities
}
