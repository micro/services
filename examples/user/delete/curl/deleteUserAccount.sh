curl "http://localhost:8080/user/Delete" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "id": "8b98acbe-0b6a-4d66-a414-5ffbf666786f"
}'