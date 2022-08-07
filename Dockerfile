FROM golang:1.18-alpine

WORKDIR /usr/src/app

COPY go.* ./
RUN go mod download

COPY . ./
RUN go build -o ./time-progression

EXPOSE 80

CMD ["./time-progression"]
