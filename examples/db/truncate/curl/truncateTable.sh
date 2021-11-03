curl "http://localhost:8080/db/Truncate" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "table": "users"
}'