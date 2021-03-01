module github.com/micro/services/sentiment

go 1.15

require (
	github.com/cdipaolo/goml v0.0.0-20190412180403-e1f51f713598 // indirect
	github.com/cdipaolo/sentiment v0.0.0-20200617002423-c697f64e7f10
	github.com/golang/protobuf v1.4.3
	github.com/micro/micro/v3 v3.1.0
	github.com/micro/services v0.0.0-00010101000000-000000000000
)

replace github.com/micro/services => ../
