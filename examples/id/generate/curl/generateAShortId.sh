curl "https://api.m3o.com/v1/id/Generate" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "type": "shortid"
}'