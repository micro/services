curl "https://api.m3o.com/v1/stock/Quote" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "symbol": "AAPL"
}'