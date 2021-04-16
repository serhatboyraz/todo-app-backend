FROM golang:1.16

WORKDIR /app

COPY . /app

RUN go get -d -v ./...

RUN go install -v ./...

EXPOSE 8080

CMD ["go", "run", "main.go"]
