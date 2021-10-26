curl "https://api.m3o.com/v1/gifs/Search" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "limit": 2,
  "query": "dogs"
}'