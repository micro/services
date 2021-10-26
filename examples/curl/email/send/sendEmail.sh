curl "https://api.m3o.com/v1/email/Send" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "from": "Awesome Dot Com",
  "subject": "Email verification",
  "textBody": "Hi there,\n\nPlease verify your email by clicking this link: $micro_verification_link"
}'