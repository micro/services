curl "http://localhost:8080/user/ResetPassword" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "code": "some-code-from-email",
  "confirmPassword": "newpass123",
  "newPassword": "newpass123"
}'