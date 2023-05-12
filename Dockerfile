FROM golang:alpine
LABEL authors="matveynematvey"

WORKDIR /main
COPY . .

RUN go mod download
RUN go build -o /main main.go

CMD ["./main"]