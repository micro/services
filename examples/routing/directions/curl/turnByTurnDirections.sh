curl "https://api.m3o.com/v1/routing/Directions" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "destination": {
    "latitude": 52.529407,
    "longitude": 13.397634
  },
  "origin": {
    "latitude": 52.517037,
    "longitude": 13.38886
  }
}'