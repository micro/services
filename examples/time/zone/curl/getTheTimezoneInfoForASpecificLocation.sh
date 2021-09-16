curl "https://api.m3o.com/v1/time/Zone" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "location": "London"
}'