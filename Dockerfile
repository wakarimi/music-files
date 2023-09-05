FROM golang:1.21 as builder
WORKDIR /go/src/app
COPY . .
RUN go mod download
WORKDIR /go/src/app/cmd/music-files
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /go/src/app/cmd/music-files/main .
CMD ["./main"]
