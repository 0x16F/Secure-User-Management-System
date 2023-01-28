FROM golang:alpine AS builder
WORKDIR /
ADD go.mod .
COPY . .
RUN go build -o app cmd/app/main.go

FROM alpine
WORKDIR /
ADD /configs /configs
ADD ./wait-for-postgres.sh ./

RUN apk update
RUN apk add postgresql-client

RUN chmod +x wait-for-postgres.sh

COPY --from=builder /app /app
CMD ["./app"]