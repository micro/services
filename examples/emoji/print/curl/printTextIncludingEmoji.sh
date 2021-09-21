curl "https://api.m3o.com/v1/emoji/Print" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "text": "let's grab a :beer:"
}'