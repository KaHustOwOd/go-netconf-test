FROM golang:1.25.7-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
#RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o go-netconf-test .

FROM alpine:3.20

WORKDIR /

COPY --from=builder /app/go-netconf-test .

ENTRYPOINT ["/go-netconf-test"]
