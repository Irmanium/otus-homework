FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/otus-homework/main.go ./
COPY ./internal ./internal
RUN CGO_ENABLED=0 GOOS=linux go build -o /otus-homework

COPY /migrations ./migrations

EXPOSE 8080

CMD ["/otus-homework"]