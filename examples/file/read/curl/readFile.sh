curl "http://localhost:8080/file/Read" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "path": "/document/text-files/file.txt",
  "project": "examples"
}'