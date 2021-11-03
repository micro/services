curl "http://localhost:8080/file/Save" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "file": {
    "content": "file content example",
    "path": "/document/text-files/file.txt",
    "project": "examples"
  }
}'