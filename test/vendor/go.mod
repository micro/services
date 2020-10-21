module github.com/micro/services/test/vendor

go 1.13

replace github.com/micro/services/test/routes => ../routes

require (
	github.com/micro/micro/v3 v3.0.0-beta.6.0.20201019094541-f64a46e81eb9
	github.com/micro/services/test/routes v0.0.0-00010101000000-000000000000
)
