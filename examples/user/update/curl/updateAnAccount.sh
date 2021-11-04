curl "http://localhost:8080/user/Update" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "email": "joeotheremail@example.com",
  "id": "usrid-1"
}'