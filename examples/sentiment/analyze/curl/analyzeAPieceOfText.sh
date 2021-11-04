curl "http://localhost:8080/sentiment/Analyze" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "text": "this is amazing"
}'