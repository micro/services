curl "http://localhost:8080/stream/CreateChannel" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "description": "The channel for all things",
  "name": "general"
}'