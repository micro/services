curl "https://api.m3o.com/v1/user/Login" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "email": "joe@example.com",
  "password": "mySecretPass123"
}'