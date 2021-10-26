curl "https://api.m3o.com/v1/location/Search" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "center": {
    "latitude": 51.511061,
    "longitude": -0.120022
  },
  "numEntities": 10,
  "radius": 100,
  "type": "bike"
}'