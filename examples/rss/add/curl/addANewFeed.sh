curl "http://localhost:8080/rss/Add" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "category": "news",
  "name": "bbc",
  "url": "http://feeds.bbci.co.uk/news/rss.xml"
}'