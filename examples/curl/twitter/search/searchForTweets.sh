curl "https://api.m3o.com/v1/twitter/Search" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "query": "cats"
}'