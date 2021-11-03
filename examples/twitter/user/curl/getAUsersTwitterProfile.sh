curl "http://localhost:8080/twitter/User" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "username": "crufter"
}'