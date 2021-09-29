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

func (e *Evchargers) loadData(r io.Reader) (int, error) {
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

func (e *Evchargers) refreshDataFromSource() {
	for {
		start := time.Now()
		logger.Infof("Refreshing data")
		// TODO replace with cron
		// wget "https://api.openchargemap.io/v3/poi/?output=json&key=<APIKEY>&maxresults=10000000" - response is json array of chargers (> 500MB).
		rsp, err := http.Get(fmt.Sprintf("https://api.openchargemap.io/v3/poi/?output=json&key=%s&maxresults=10000000", e.conf.OCMKey))
		if err != nil {
			logger.Errorf("Error refreshing data %s", err)
		} else {
			logger.Infof("Loading data")
			c, err := e.loadData(rsp.Body)
			if err != nil {
				logger.Errorf("Error loading data %s", err)
			} else {
				logger.Infof("Updated %v. Took %s", c, time.Since(start))
			}
			rsp.Body.Close()
		}

		time.Sleep(24 * time.Hour)
	}
}