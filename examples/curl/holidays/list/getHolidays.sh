curl "https://api.m3o.com/v1/holidays/List" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "country_code": "GB",
  "year": 2022
}'