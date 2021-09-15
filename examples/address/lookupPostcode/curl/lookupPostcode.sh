curl "https://api.m3o.com/v1/address/LookupPostcode" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "postcode": "SW1A 2AA"
}'