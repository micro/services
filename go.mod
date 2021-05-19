module github.com/micro/services

go 1.14

require (
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/disintegration/imaging v1.6.2
	github.com/getkin/kin-openapi v0.26.0
	github.com/gojuno/go.osrm v0.1.1-0.20200217151037-435fc3e1d3d4
	github.com/golang/protobuf v1.5.1
	github.com/google/uuid v1.1.2
	github.com/hailocab/go-geoindex v0.0.0-20160127134810-64631bfe9711
	github.com/hashicorp/golang-lru v0.5.3
	github.com/lib/pq v1.9.0 // indirect
	github.com/micro/dev v0.0.0-20201117163752-d3cfc9788dfa
	github.com/micro/micro/v3 v3.2.2-0.20210514120220-1ee39904d3dd
	github.com/miekg/dns v1.1.31 // indirect
	github.com/paulmach/go.geo v0.0.0-20180829195134-22b514266d33
	github.com/pquerna/otp v1.3.0
	github.com/stoewer/go-strcase v1.2.0
	github.com/stretchr/testify v1.7.0
	github.com/teris-io/shortid v0.0.0-20171029131806-771a37caa5cf
	github.com/ulikunitz/xz v0.5.8 // indirect
	golang.org/x/crypto v0.0.0-20210220033148-5ea612d1eb83
	golang.org/x/net v0.0.0-20201021035429-f5854403a974
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43 // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a // indirect
	golang.org/x/sys v0.0.0-20210225134936-a50acf3fe073 // indirect
	google.golang.org/genproto v0.0.0-20201001141541-efaab9d3c4f7 // indirect
	google.golang.org/grpc v1.32.0 // indirect
	google.golang.org/protobuf v1.26.0
	googlemaps.github.io/maps v1.3.1
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
	gorm.io/driver/postgres v1.0.6
	gorm.io/gorm v1.20.9
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
