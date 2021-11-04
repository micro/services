curl "http://localhost:8080/rss/Feed" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "name": "bbc"
}'