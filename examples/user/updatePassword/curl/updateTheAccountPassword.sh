curl "http://localhost:8080/user/UpdatePassword" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "confirmPassword": "myEvenMoreSecretPass123",
  "id": "usrid-1",
  "newPassword": "myEvenMoreSecretPass123",
  "oldPassword": "mySecretPass123"
}'