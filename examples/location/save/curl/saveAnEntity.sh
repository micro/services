curl "https://api.m3o.com/v1/location/Save" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "entity": {
    "id": "1",
    "location": {
      "latitude": 51.511061,
      "longitude": -0.120022,
      "timestamp": "1622802761"
    },
    "type": "bike"
  }
}'