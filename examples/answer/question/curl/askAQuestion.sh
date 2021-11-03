curl "http://localhost:8080/answer/Question" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "query": "microsoft"
}'