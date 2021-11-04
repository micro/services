curl "http://localhost:8080/prayer/Times" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "location": "london"
}'