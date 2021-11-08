curl "http://localhost:8080/helloworld/Stream" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "messages": 10,
  "name": "John"
}'