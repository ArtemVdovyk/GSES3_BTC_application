FROM golang:1.20.4-alpine

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY GSES3_BTC_application/*.go GSES3_BTC_application/.env ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /gses3_btc_application .

EXPOSE 8080

ENTRYPOINT ["/gses3_btc_application"]
