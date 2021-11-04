curl "http://localhost:8080/weather/Forecast" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "days": 2,
  "location": "London"
}'