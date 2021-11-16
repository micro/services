curl "http://localhost:8080/user/SendPasswordResetEmail" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "email": "joe@example.com",
  "fromName": "Awesome Dot Com",
  "subject": "Password reset",
  "textContent": "Hi there,\n click here to reset your password: myapp.com/reset/code?=$code"
}'