FROM golang:1-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o rr ./cmd/rr

FROM alpine:latest
WORKDIR /request-review
COPY --from=builder /request-review/rr .
ENTRYPOINT [ "/request-review/rr" ]
