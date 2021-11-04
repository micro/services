curl "http://localhost:8080/id/Generate" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "type": "shortid"
}'