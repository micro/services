curl "https://api.m3o.com/v1/emoji/Find" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "alias": ":beer:"
}'