curl "http://localhost:8080/db/Count" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "table": "users"
}'