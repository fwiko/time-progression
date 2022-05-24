FROM golang:1.18-alpine

workdir /usr/src/app

copy go.mod .
copy go.sum .
RUN go mod download

copy . .
RUN go build -o ./time-progression

EXPOSE 80

CMD ["./time-progression"]