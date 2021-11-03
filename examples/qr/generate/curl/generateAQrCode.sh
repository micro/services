curl "http://localhost:8080/qr/Generate" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "size": 300,
  "text": "https://m3o.com/qr"
}'