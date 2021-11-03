curl "http://localhost:8080/emoji/Flag" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "alias": "GB"
}'