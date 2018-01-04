FROM golang:1.8.3-alpine as builder

WORKDIR /go/src/packagetree
COPY . .

RUN go test -cover ./... \
    && go build -o reposerver cmd/reposerver/reposerver.go

FROM alpine:3.6

EXPOSE 8080

WORKDIR /run/
COPY --from=builder /go/src/packagetree/reposerver .

RUN apk --no-cache add ca-certificates

USER nobody
ENTRYPOINT ["./reposerver"]