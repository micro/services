curl "http://localhost:8080/thumbnail/Screenshot" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "height": 600,
  "url": "https://m3o.com",
  "width": 600
}'