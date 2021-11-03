curl "https://api.m3o.com/v1/mq/Subscribe" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "topic": "events"
}'