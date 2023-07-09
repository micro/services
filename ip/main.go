package main

import (
	"io"
	"os"
	"path"
	"strings"

	"github.com/micro/services/ip/handler"
	pb "github.com/micro/services/ip/proto"

	"github.com/oschwald/geoip2-golang"
	"micro.dev/v4/service"
	"micro.dev/v4/service/config"
	"micro.dev/v4/service/logger"
	"micro.dev/v4/service/store"
)

// loadFile from the blob store
func loadFile(p string) (string, error) {
	name := path.Base(p)

	f, err := os.Create("./" + name)
	if err != nil {
		return "", err
	}
	defer f.Close()

	reader, err := store.DefaultBlobStore.Read(p)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(f, reader)
	return "./" + name, err
}

func main() {
	// Create service
	srv := service.New(
		service.Name("ip"),
		service.Version("latest"),
	)

	// get the ip city database
	v, err := config.Get("ip.city.database")
	if err != nil {
		logger.Fatalf("failed to get config: %v", err)
	}
	path := v.String("./GeoLite2-City.mmdb")

	// load from blob store if specified
	if strings.HasPrefix(path, "blob://") {
		f, err := loadFile(strings.TrimPrefix(path, "blob://"))
		if err != nil {
			logger.Fatal("failed to load db: %v", err)
		}

		path = f
	}

	// load the ip city database
	cr, err := geoip2.Open(path)
	if err != nil {
		logger.Fatalf("failed to open ip db: %v", err)
	}

	// get the asn database
	v, err = config.Get("ip.asn.database")
	if err != nil {
		logger.Fatalf("failed to get config: %v", err)
	}
	path = v.String("./GeoLite2-ASN.mmdb")

	// load from blob store if specified
	if strings.HasPrefix(path, "blob://") {
		f, err := loadFile(strings.TrimPrefix(path, "blob://"))
		if err != nil {
			logger.Fatal("failed to load db: %v", err)
		}

		path = f
	}

	ar, err := geoip2.Open(path)
	if err != nil {
		logger.Fatalf("failed to open ip db: %v", err)
	}

	// Register handler
	pb.RegisterIpHandler(srv.Server(), &handler.Ip{CityReader: cr, ASNReader: ar})

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
