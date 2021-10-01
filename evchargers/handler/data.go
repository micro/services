package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/micro/micro/v3/service/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (e *Evchargers) loadPOIData(r io.Reader) (int, error) {
	logger.Infof("Loading reference data")
	dec := json.NewDecoder(r)
	t, err := dec.Token()
	if err != nil {
		return 0, err
	}
	d, ok := t.(json.Delim)
	if !ok || d.String() != "[" {
		return 0, fmt.Errorf("unexpected token %v %+v", ok, t)
	}
	ctx := context.Background()
	count := 0
	for dec.More() {
		// process each item in json array and insert into mongodb
		var p Poi
		if err := dec.Decode(&p); err != nil {
			return 0, fmt.Errorf("error unmarshalling charger %s", err)
		}
		if len(p.SpatialPosition.Type) == 0 {
			// blank so reconstruct
			p.SpatialPosition.Type = "Point"
			// long, lat not lat, long
			p.SpatialPosition.Coordinates = []float64{p.Address.Longitude, p.Address.Latitude}
		}

		t := true
		_, err := e.mdb.Database("ocm").Collection("poi").ReplaceOne(ctx, bson.D{bson.E{"ID", p.ID}}, p, &options.ReplaceOptions{Upsert: &t})
		if err != nil {
			return 0, err
		}
		count++
	}
	return count, nil

}

func (e *Evchargers) loadRefData(r io.Reader) error {
	dec := json.NewDecoder(r)
	var rd ReferenceData
	if err := dec.Decode(&rd); err != nil {
		return err
	}
	ctx := context.Background()
	t := true
	_, err := e.mdb.Database("ocm").Collection("reference").ReplaceOne(ctx, bson.D{bson.E{"_id", 1}}, rd, &options.ReplaceOptions{Upsert: &t})
	if err != nil {
		return err
	}

	return nil
}

func (e *Evchargers) refreshDataFromSource() {
	start := time.Now()
	logger.Infof("Refreshing data")
	logger.Infof("Retrieving poi data")
	rsp, err := http.Get(fmt.Sprintf("https://api.openchargemap.io/v3/poi/?output=json&key=%s&maxresults=10000000", e.conf.OCMKey))
	if err != nil {
		logger.Errorf("Error refreshing data %s", err)
		return
	}
	defer rsp.Body.Close()
	c, err := e.loadPOIData(rsp.Body)
	if err != nil {
		logger.Errorf("Error loading data %s", err)
		return
	}
	logger.Infof("Updated %v items of POI data. Took %s", c, time.Since(start))

	start = time.Now()
	logger.Infof("Retrieving ref data")
	rsp2, err := http.Get(fmt.Sprintf("https://api.openchargemap.io/v3/referencedata/?output=json&key=%s", e.conf.OCMKey))
	if err != nil {
		logger.Errorf("Error refreshing reference data %s", err)
		return
	}
	defer rsp2.Body.Close()
	if err := e.loadRefData(rsp2.Body); err != nil {
		logger.Errorf("Error loading reference data %s", err)
		return
	}
	logger.Infof("Updated ref data. Took %s", time.Since(start))

}
