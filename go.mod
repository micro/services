module github.com/micro/services

go 1.14

require (
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/PuerkitoBio/goquery v1.6.1
	github.com/SlyMarbo/rss v1.0.1
	github.com/cdipaolo/goml v0.0.0-20190412180403-e1f51f713598 // indirect
	github.com/cdipaolo/sentiment v0.0.0-20200617002423-c697f64e7f10
	github.com/crufter/lexer v0.0.0-20120907053443-23fe8c7add01
	github.com/disintegration/imaging v1.6.2
	github.com/getkin/kin-openapi v0.26.0
	github.com/gojuno/go.osrm v0.1.1-0.20200217151037-435fc3e1d3d4
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/golang/protobuf v1.5.1
	github.com/google/uuid v1.1.2
	github.com/hailocab/go-geoindex v0.0.0-20160127134810-64631bfe9711
	github.com/hashicorp/golang-lru v0.5.3
	github.com/jackc/pgx/v4 v4.10.1
	github.com/lib/pq v1.9.0 // indirect
	github.com/micro/dev v0.0.0-20201117163752-d3cfc9788dfa
	github.com/micro/micro/v3 v3.2.2-0.20210607154842-ec8964031a93
	github.com/miekg/dns v1.1.31 // indirect
	github.com/onsi/ginkgo v1.15.0 // indirect
	github.com/oschwald/geoip2-golang v1.5.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/paulmach/go.geo v0.0.0-20180829195134-22b514266d33
	github.com/pquerna/otp v1.3.0
	github.com/stoewer/go-strcase v1.2.0
	github.com/stretchr/testify v1.7.0
	github.com/teris-io/shortid v0.0.0-20171029131806-771a37caa5cf
	go.opencensus.io v0.22.4 // indirect
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2
	golang.org/x/net v0.0.0-20210226172049-e18ecbb05110
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a // indirect
	golang.org/x/sys v0.0.0-20210225134936-a50acf3fe073 // indirect
	google.golang.org/genproto v0.0.0-20201001141541-efaab9d3c4f7 // indirect
	google.golang.org/grpc v1.32.0 // indirect
	google.golang.org/protobuf v1.26.0
	googlemaps.github.io/maps v1.3.1
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	gorm.io/datatypes v1.0.1
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.21.10
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
