FROM golang:1.20.5-alpine3.18 as builder

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

FROM alpine:3.18

RUN apk add bash

COPY --from=builder /go/bin/migrate /usr/bin/go-migrate

RUN chmod +x /usr/bin/go-migrate

COPY ./docker/migrate.sh /usr/bin/migrate

RUN chmod +x /usr/bin/migrate

WORKDIR /root

COPY ./backend/sql/migrations ./migrations

CMD [ "bash" ]
