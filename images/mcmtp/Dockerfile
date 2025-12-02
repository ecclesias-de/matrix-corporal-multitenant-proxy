FROM golang:1.24.3 as builder

ARG CGO_ENABLED=0
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build cmd/corporal-multitenant-proxy.go

FROM scratch
COPY --from=builder /app/corporal-multitenant-proxy /mcmtp
CMD ["/mcmtp" "-config" "/etc/mcmtp/config.json"]