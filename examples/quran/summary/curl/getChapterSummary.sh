curl "https://api.m3o.com/v1/quran/Summary" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "chapter": 1
}'