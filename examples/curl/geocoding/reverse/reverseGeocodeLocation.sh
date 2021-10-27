curl "https://api.m3o.com/v1/geocoding/Reverse" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "latitude": 51.5123064,
  "longitude": -0.1216235
}'