curl "https://api.m3o.com/v1/stock/OrderBook" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "date": "2020-10-01",
  "end": "2020-10-01T11:00:00Z",
  "limit": 3,
  "start": "2020-10-01T10:00:00Z",
  "stock": "AAPL"
}'