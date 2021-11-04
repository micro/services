curl "http://localhost:8080/helloworld/Call" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "name": "John"
}'