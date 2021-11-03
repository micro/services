curl "http://localhost:8080/sunnah/Books" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "collection": "bukhari"
}'