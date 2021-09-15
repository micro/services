curl "https://api.m3o.com/v1/currency/Convert" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "amount": 10,
  "from": "USD",
  "to": "GBP"
}'