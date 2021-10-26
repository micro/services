curl "https://api.m3o.com/v1/weather/Forecast" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "days": 2,
  "location": "London"
}'