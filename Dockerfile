FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /twitter-go-api

EXPOSE 8080

CMD [ "/twitter-go-api" ]