FROM node:lts-alpine3.16 AS node_builder

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

FROM nginx:alpine

RUN rm /etc/nginx/conf.d/default.conf

COPY ./docker/nginx-static-setup/nginx.conf /etc/nginx/conf.d/webapp.conf

COPY --from=node_builder /build/frontend/dist /var/www/www-root

RUN mkdir -p /var/www/www-root/auth/signin
RUN mkdir -p /var/www/www-root/auth/signup

RUN cp /var/www/www-root/index.html /var/www/www-root/auth/signup/
RUN cp /var/www/www-root/index.html /var/www/www-root/auth/signin/

RUN chown -R nginx:nginx /var/www/www-root
