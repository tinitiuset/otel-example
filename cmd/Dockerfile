FROM alpine:latest

WORKDIR /app
COPY wave-generator .

ENV OTEL_COLLECTOR_ENDPOINT=localhost:4317

CMD ["./wave-generator"] 