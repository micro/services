curl "http://localhost:8080/thumbnail/Screenshot" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "height": 600,
  "url": "https://google.com",
  "width": 600
}'