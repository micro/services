curl "https://api.m3o.com/v1/google/Search" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "query": "how to make donuts"
}'