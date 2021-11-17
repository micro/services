curl "http://localhost:8080/user/Update" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "email": "joe+2@example.com",
  "id": "user-1"
}'