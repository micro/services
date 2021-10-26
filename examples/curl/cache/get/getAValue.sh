curl "https://api.m3o.com/v1/cache/Get" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "key": "foo"
}'