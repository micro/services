curl "https://api.m3o.com/v1/event/Read" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "topic": "user"
}'