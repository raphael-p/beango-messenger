FROM golang:1.21.3

WORKDIR /usr/src/app/

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod ./
RUN go mod tidy
