curl "https://api.m3o.com/v1/function/Deploy" \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $MICRO_API_TOKEN" \
-d '{
  "entrypoint": "helloworld",
  "name": "my-first-func",
  "project": "tests",
  "repo": "github.com/m3o/nodejs-function-example",
  "runtime": "nodejs14"
}'