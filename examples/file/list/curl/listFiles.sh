curl "http://localhost:8080/file/List" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "project": "examples"
}'