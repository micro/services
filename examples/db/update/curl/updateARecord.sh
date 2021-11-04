curl "http://localhost:8080/db/Update" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "record": {
    "age": 43,
    "id": "1"
  },
  "table": "users"
}'