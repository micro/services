module github.com/micro/services

go 1.14

require (
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.2
	github.com/gosimple/slug v1.9.0
	github.com/m3o/services v0.0.0-20201013090801-c9adc79659de
	github.com/micro/go-micro/v3 v3.0.0-beta.3.0.20201014103235-c49fef49b8af
	github.com/micro/micro/plugin/etcd/v3 v3.0.0-20201014133532-d4ad235f4987 // indirect
	github.com/micro/micro/v3 v3.0.0-beta.6.0.20201014125737-097dd92a1b29
	github.com/miekg/dns v1.1.31 // indirect
	github.com/ulikunitz/xz v0.5.8 // indirect
	golang.org/x/crypto v0.0.0-20201002094018-c90954cbb977 // indirect
	golang.org/x/net v0.0.0-20200930145003-4acb6c075d10 // indirect
	golang.org/x/oauth2 v0.0.0-20200902213428-5d25da1a8d43 // indirect
	golang.org/x/sys v0.0.0-20200930185726-fdedc70b468f // indirect
	google.golang.org/genproto v0.0.0-20201013134114-7f9ee70cb474 // indirect
	google.golang.org/grpc v1.32.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
