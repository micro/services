curl "http://localhost:8080/user/Login" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "email": "joe@example.com",
  "password": "mySecretPass123"
}'