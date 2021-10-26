curl "https://api.m3o.com/v1/notes/Update" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "note": {
    "id": "63c0cdf8-2121-11ec-a881-0242e36f037a",
    "text": "Updated note text",
    "title": "Update Note"
  }
}'