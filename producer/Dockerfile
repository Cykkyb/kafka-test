FROM golang:latest

RUN go version
ENV GOPATH=/

COPY . .

RUN go build  -o main-app ./cmd/main.go

CMD  ./main-app