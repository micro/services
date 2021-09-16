curl "https://api.m3o.com/v1/stream/Publish" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "message": {
    "id": "1",
    "type": "signup",
    "user": "john"
  },
  "topic": "events"
}'