FROM golang:latest

RUN go version
ENV GOPATH=/

COPY . .
COPY migrations /migrations

RUN go build  -o main-app ./cmd/main.go
RUN go build  -o migrate ./cmd/migrate/main.go

CMD ./migrate && ./main-app