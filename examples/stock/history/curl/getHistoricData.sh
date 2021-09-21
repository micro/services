curl "https://api.m3o.com/v1/stock/History" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "date": "2020-10-01",
  "stock": "AAPL"
}'