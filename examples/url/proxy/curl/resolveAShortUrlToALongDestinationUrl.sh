curl "http://localhost:8080/url/Proxy" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "shortURL": "https://m3o.one/u/ck6SGVkYp"
}'