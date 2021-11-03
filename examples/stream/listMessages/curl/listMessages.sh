curl "https://api.m3o.com/v1/stream/ListMessages" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "channel": "general"
}'