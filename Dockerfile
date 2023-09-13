ARG GOLANG_VERSION=1.21
ARG ALPINE_VERSION=3.18

###
# BUILD
###
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} as builder

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/test-webserver ./...


###
# APP
###
FROM alpine:${ALPINE_VERSION}

COPY --from=builder /usr/local/bin/test-webserver /usr/local/bin/test-webserver

CMD ["/usr/local/bin/test-webserver"]