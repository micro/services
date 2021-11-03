curl "http://localhost:8080/notes/Create" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "text": "This is my note",
  "title": "New Note"
}'