curl "http://localhost:8080/user/UpdatePassword" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "confirmPassword": "Password2",
  "id": "user-1",
  "newPassword": "Password2",
  "oldPassword": "Password1"
}'