curl "https://api.m3o.com/v1/twitter/User" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "username": "crufter"
}'