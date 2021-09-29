curl "https://api.m3o.com/v1/notes/Create" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "text": "This is my note",
  "title": "New Note"
}'