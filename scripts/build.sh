docker run -it --rm \
    -e GOOS=linux \
    -e GOARCH=arm \
    -e CGO_ENABLED=0 \
    -v "$(pwd)":/work \
    -w /work \
    golang:1.17 \
    go build -ldflags '-w -extldflags "-static"' -o door .
