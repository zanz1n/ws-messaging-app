FROM nginx:alpine

RUN rm /etc/nginx/conf.d/default.conf

COPY ./docker/nginx.conf /etc/nginx/conf.d/webapp.conf
