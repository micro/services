curl "http://localhost:8080/db/Create" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "record": {
    "age": 42,
    "id": "1",
    "isActive": true,
    "name": "Jane"
  },
  "table": "users"
}'