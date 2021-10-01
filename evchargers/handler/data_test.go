package handler

import (
	"context"
	"io"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestDataLoad(t *testing.T) {
	t.SkipNow()
	s := Evchargers{}
	opts := []*options.ClientOptions{options.Client().ApplyURI("mongodb://127.0.0.1:27017/ocm")}
	var err error
	s.mdb, err = mongo.Connect(context.Background(), opts...)
	if err != nil {
		t.Fatalf("Error connecting to mongo %s", err)
	}
	var r io.Reader
	r, err = os.Open("test.json")
	if err != nil {
		t.Fatalf("Error opening test data %s", err)
	}
	c, err := s.loadPOIData(r)
	if err != nil {
		t.Fatalf("Err loading data %s", err)
	}
	if c != 2 {
		t.Errorf("Incorrect number of records %d", c)
	}
}

func TestRefDataLoad(t *testing.T) {
	t.SkipNow()
	s := Evchargers{}
	opts := []*options.ClientOptions{options.Client().ApplyURI("mongodb://127.0.0.1:27017/ocm")}
	var err error
	s.mdb, err = mongo.Connect(context.Background(), opts...)
	if err != nil {
		t.Fatalf("Error connecting to mongo %s", err)
	}
	var r io.Reader
	r, err = os.Open("test-reference.json")
	if err != nil {
		t.Fatalf("Error opening test data %s", err)
	}
	err = s.loadRefData(r)
	if err != nil {
		t.Fatalf("Err loading data %s", err)
	}

}
