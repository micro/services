curl "http://localhost:8080/emoji/Print" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "text": "let's grab a :beer:"
}'