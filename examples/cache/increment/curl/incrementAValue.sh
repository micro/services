curl "http://localhost:8080/cache/Increment" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "key": "counter",
  "value": 2
}'