curl "https://api.m3o.com/v1/crypto/News" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "symbol": "BTCUSD"
}'