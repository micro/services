curl "https://api.m3o.com/v1/user/Read" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "id": "usrid-1"
}'