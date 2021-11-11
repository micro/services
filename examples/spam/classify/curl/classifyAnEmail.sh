curl "http://localhost:8080/spam/Classify" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "email_body": "Hi there,\n\nWelcome to M3O.\n\nThanks\nM3O team",
  "from": "noreply@m3o.com",
  "subject": "Welcome",
  "to": "hello@example.com"
}'