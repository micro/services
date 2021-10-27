curl "https://api.m3o.com/v1/youtube/Search" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "query": "donuts"
}'