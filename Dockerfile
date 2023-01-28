FROM golang:alpine AS builder
WORKDIR /
ADD go.mod .
COPY . .
RUN go build -o app cmd/app/main.go

FROM alpine
WORKDIR /

ADD /configs /configs

RUN apk update
RUN apk add postgresql-client

COPY --from=builder /app /app
CMD ["./app"]