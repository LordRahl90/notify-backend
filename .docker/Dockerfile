FROM golang:alpine AS builder
ADD . /app/backend
WORKDIR /app/backend
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o notifier
# RUN go build -o notifier

# final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/backend/ ./

RUN chmod +x ./notifier
ENTRYPOINT ["./notifier"]
CMD ["-env=docker"]
EXPOSE 5500