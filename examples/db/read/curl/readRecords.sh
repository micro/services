curl "https://api.m3o.com/v1/db/Read" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "query": "age == 43",
  "table": "users"
}'