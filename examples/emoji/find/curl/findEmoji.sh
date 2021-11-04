curl "http://localhost:8080/emoji/Find" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "alias": ":beer:"
}'