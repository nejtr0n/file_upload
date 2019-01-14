FROM golang:1.11.4 as builder
RUN mkdir /build
ADD . /build/storage
WORKDIR /build/storage
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main pkg/web.go

FROM alpine:latest
RUN adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/storage/main /app/main
COPY --from=builder /build/storage/ui/http/templates /app/storage/ui/http/templates
WORKDIR /app
CMD ["./main"]