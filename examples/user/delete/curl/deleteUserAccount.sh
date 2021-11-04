curl "http://localhost:8080/user/Delete" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "id": "fdf34f34f34-f34f34-f43f43f34-f4f34f"
}'