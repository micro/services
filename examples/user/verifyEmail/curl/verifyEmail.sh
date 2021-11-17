curl "http://localhost:8080/user/VerifyEmail" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "email": "joe@example.com",
  "token": "012345"
}'