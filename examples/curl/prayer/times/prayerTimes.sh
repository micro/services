curl "https://api.m3o.com/v1/prayer/Times" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "location": "london"
}'