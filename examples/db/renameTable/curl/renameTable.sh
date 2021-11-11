curl "http://localhost:8080/db/RenameTable" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "from": "events",
  "to": "events_backup"
}'