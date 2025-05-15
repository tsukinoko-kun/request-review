FROM golang:1-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o rr ./cmd/request-review

FROM alpine:latest
WORKDIR /request-review
COPY --from=builder /app/rr .
RUN apk add --no-cache git openssh-client
ENV GIT_TERMINAL_PROMPT=0
ENTRYPOINT [ "/request-review/rr" ]
