curl "http://localhost:8080/user/Read" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "username": "usrname-1"
}'