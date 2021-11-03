curl "http://localhost:8080/cache/Get" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "key": "foo"
}'