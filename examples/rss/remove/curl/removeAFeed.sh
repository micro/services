curl "https://api.m3o.com/v1/rss/Remove" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "name": "bbc"
}'