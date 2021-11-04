curl "http://localhost:8080/user/Logout" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "sessionId": "sds34s34s34-s34s34-s43s43s34-s4s34s"
}'