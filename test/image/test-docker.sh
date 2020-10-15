cp services/tests/image/.dockerignore .
pushd micro
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build
popd
DOCKER_BUILDKIT=1 docker build -t micro -f services/tests/image/Dockerfile .
