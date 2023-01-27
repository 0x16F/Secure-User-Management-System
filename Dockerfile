FROM golang:alpine AS builder
WORKDIR /
ADD go.mod .
COPY . .
RUN go build -o app cmd/app/main.go
FROM alpine
WORKDIR /
ADD /configs /configs
COPY --from=builder /app /app
CMD ["./app"]