FROM node:lts-alpine3.18 AS node_builder

WORKDIR /build

RUN npm i -g pnpm

COPY ./frontend/package.json /build/frontend/

COPY ./package.json ./pnpm-lock.yaml ./pnpm-workspace.yaml /build/

RUN pnpm install --frozen-lockfile

COPY ./frontend/index.html \
    ./frontend/vite.config.ts \
    ./frontend/env-settings.json \
    ./frontend/tsconfig.json \
    ./frontend/tsconfig.node.json \
    /build/frontend/

COPY ./frontend/public /build/frontend/public

COPY ./frontend/src /build/frontend/src

RUN pnpm build:client

FROM golang:1.20.5-alpine3.18 as builder

WORKDIR /build

RUN apk add make

COPY ./backend /build/backend

COPY ./go.work ./go.work.sum /build/

WORKDIR /build/backend

RUN make build-min

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /build/backend/dist/stable.bin /app/bin
COPY --from=node_builder /build/frontend/dist /app/static

CMD [ "/app/bin" ]
