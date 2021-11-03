curl "http://localhost:8080/crypto/Quote" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "symbol": "BTCUSD"
}'