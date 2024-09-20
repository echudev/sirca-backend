FROM golang:1.23.0

WORKDIR /usr/src/app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /sirca-backend ./cmd/server

EXPOSE 8080

CMD [ "/sirca-backend" ]
