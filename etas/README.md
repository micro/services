# ETAs Service

This is the ETAs service. It provides ETAs for single-pickup, multi-dropoff routes. It takes into account time and traffic.

Current limitations:
• Only supports "Driving" (not walking, cycling)
• Does not optimize route

## Usage

The ETA service depends on the Google Maps API. Ensure you set the "google.maps.apikey" config value to your API key.

```
micro config set google.maps.apikey YOUR_API_KEY
```

Once set, run the service using `micro run github.com/micro/services/etas`.

```bash
$ micro call etas ETAs.Calculate $(cat example-req.json)
{
	"points": {
		"brentwood-station": {
			"estimated_arrival_time": "2020-12-15T11:01:29.429947Z",
			"estimated_departure_time": "2020-12-15T11:01:29.429947Z"
		},
		"nandos": {
			"estimated_arrival_time": "2020-12-15T10:54:38.429947Z",
			"estimated_departure_time": "2020-12-15T10:54:38.429947Z"
		},
		"shenfield-station": {
			"estimated_arrival_time": "2020-12-15T10:48:34.429947Z",
			"estimated_departure_time": "2020-12-15T10:48:34.429947Z"
		}
	}
}
```
