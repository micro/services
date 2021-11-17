curl "http://localhost:8080/user/Logout" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "sessionId": "df91a612-5b24-4634-99ff-240220ab8f55"
}'