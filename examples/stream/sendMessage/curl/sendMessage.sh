curl "http://localhost:8080/stream/SendMessage" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "channel": "general",
  "text": "Hey checkout this tweet https://twitter.com/m3oservices/status/1455291054295498752"
}'