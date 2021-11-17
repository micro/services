curl "http://localhost:8080/user/ResetPassword" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "code": "012345",
  "confirmPassword": "NewPassword1",
  "email": "joe@example.com",
  "newPassword": "NewPassword1"
}'