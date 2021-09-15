curl "https://api.m3o.com/v1/helloworld/Call" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "name": "John"
}'