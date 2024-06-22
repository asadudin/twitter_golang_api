FROM golang:1.18-alpine

# Install git for go module downloads
RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /twitter-go-api

EXPOSE 8080

CMD [ "/twitter-go-api" ]
