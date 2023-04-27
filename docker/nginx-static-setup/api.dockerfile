FROM golang:1.20.3-alpine3.16 as builder

WORKDIR /build

RUN apk add make

COPY ./backend /build/backend

COPY ./go.work ./go.work.sum /build/

WORKDIR /build/backend

RUN make build-min

FROM alpine:3.16

COPY --from=builder /build/backend/dist/stable.bin /app

CMD [ "/app" ]
