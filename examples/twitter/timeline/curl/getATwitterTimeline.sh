curl "http://localhost:8080/twitter/Timeline" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "limit": 1,
  "username": "m3oservices"
}'