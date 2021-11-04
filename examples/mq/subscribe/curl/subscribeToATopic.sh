curl "http://localhost:8080/mq/Subscribe" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "topic": "events"
}'