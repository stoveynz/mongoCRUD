FROM golang:alpine AS builder
RUN apk update && apk add git
WORKDIR /build
ENV GOSUMDB=off
ENV GOPROXY=direct


COPY go.mod ./
RUN go mod download

COPY *.go ./
RUN go get 
RUN go build -tags netgo -a -v -o test-api .

FROM alpine:latest
WORKDIR /opt/server
ARG ARCH

COPY --from=builder /build/test-api /usr/local/bin/ 
RUN chmod +x /usr/local/bin/test-api

EXPOSE 8000
CMD ["test-api"]