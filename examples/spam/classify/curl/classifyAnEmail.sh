curl "http://localhost:8080/spam/Classify" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "from": "noreply@m3o.com",
  "subject": "Welcome",
  "text_body": "Hi there,\n\nWelcome to M3O.\n\nThanks\nM3O team",
  "to": "hello@example.com"
}'