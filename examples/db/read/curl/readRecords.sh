curl "http://localhost:8080/db/Read" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "query": "age == 43",
  "table": "users"
}'