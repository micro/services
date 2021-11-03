curl "http://localhost:8080/notes/Events" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "id": "63c0cdf8-2121-11ec-a881-0242e36f037a"
}'