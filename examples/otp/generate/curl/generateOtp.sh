curl "http://localhost:8080/otp/Generate" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "id": "asim@example.com"
}'