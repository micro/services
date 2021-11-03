curl "http://localhost:8080/currency/Rates" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "code": "USD"
}'