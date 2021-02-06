FROM golang:latest
RUN mkdir -p /go/src/app
WORKDIR /go/src/app
COPY . /go/src/app
COPY go.mod go.sum ./

RUN go mod download

COPY . .

