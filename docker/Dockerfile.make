FROM golang:latest AS builder
MAINTAINER Jake Sylvestre <jake@entropy.rocks>
WORKDIR /updater

RUN apt install gcc

ADD ./go.mod ./
ADD ./go.sum ./
ADD ./main.go ./
COPY ./pkg ./pkg
COPY vendor ./vendor

RUN CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -mod vendor -o updater

FROM alpine:latest
MAINTAINER Jake Sylvestre <jake@entropy.rocks>
WORKDIR /app
COPY --from=builder /updater/updater /app/updater
RUN apk add --no-cache libc6-compat

CMD [ "/app/updater" ]