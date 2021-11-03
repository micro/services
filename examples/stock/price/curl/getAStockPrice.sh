curl "http://localhost:8080/stock/Price" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "symbol": "AAPL"
}'