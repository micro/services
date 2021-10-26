curl "https://api.m3o.com/v1/function/Describe" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "name": "my-first-func",
  "project": "tests"
}'