curl "http://localhost:8080/postcode/Lookup" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "postcode": "SW1A 2AA"
}'