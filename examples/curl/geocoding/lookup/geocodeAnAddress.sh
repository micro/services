curl "https://api.m3o.com/v1/geocoding/Lookup" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "address": "10 russell st",
  "city": "london",
  "country": "uk",
  "postcode": "wc2b"
}'