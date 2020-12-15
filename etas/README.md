# ETAs Service

This is the ETAs service. It provides ETAs for single-pickup, multi-dropoff routes. It takes into account time and traffic.

Current limitations:
• Only supports "Driving" (not walking, cycling)
• Does not optimize route

## Usage

There is one required config value: `google.maps.apikey`. Once you have set this config value, run the service using `micro run`.

```bash
micro@Bens-MBP-3 etas % micro call etas ETAs.Calculate $(cat example-req.json)
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
