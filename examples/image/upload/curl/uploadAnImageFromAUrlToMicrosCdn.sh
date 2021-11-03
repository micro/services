curl "http://localhost:8080/image/Upload" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "name": "cat.jpeg",
  "url": "somewebsite.com/cat.png"
}'