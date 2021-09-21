curl "https://api.m3o.com/v1/otp/Validate" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "code": "656211",
  "id": "asim@example.com"
}'