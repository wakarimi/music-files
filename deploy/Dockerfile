# Build the binary
FROM golang:1.22 AS build

RUN useradd -u 10001 gopher

ENV NAME "wakarimi-files"
RUN mkdir /${NAME}
WORKDIR /opt/${NAME}

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ./bin/${NAME} ./cmd/main.go

# Run the binary
FROM scratch

ENV NAME "wakarimi-files"

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd

USER gopher

COPY --from=build /opt/${NAME}/bin/${NAME} /${NAME}
COPY --from=build /opt/${NAME}/migration /migration

CMD ["./wakarimi-files"]