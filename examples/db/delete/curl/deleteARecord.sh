curl "http://localhost:8080/db/Delete" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "id": "1",
  "table": "users"
}'