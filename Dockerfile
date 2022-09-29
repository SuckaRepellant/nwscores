FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o /usr/bin/nwscores

EXPOSE 8000

CMD ["/usr/bin/nwscores"]
