curl "https://api.m3o.com/v1/vehicle/Lookup" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "registration": "LC60OTA"
}'