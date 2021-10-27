curl "https://api.m3o.com/v1/evchargers/Search" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "distance": 2000,
  "levels": [
    "3"
  ],
  "location": {
    "latitude": 51.53336351319885,
    "longitude": -0.0252
  },
  "max_results": 2
}'