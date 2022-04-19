#build stage
FROM golang:1.18-buster AS builder
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go build -o /go/bin/app

#final stage
FROM gcr.io/distroless/base-debian10
COPY --from=builder /go/bin/app /app
ENTRYPOINT ["/app"]