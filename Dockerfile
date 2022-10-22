FROM golang:1.19-alpine as builder
RUN apk --update --no-cache add g++

WORKDIR /build

RUN go mod init
RUN go mod tidy

COPY . .

RUN go build -a --ldflags '-linkmode external -extldflags "-static"' .

FROM alpine
RUN apk --update --no-cache add ca-certificates \
    && addgroup -S logsrv && adduser -S -g logsrv logsrv
USER logsrv

ENV LOGSRV_HOST=0.0.0.0 LOGSRV_PORT=8080
ENTRYPOINT ["/logsrv"]
EXPOSE 8080

COPY --from=builder /build/logsrv /