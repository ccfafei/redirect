FROM node:18

WORKDIR /redirect_echart/
COPY . .

RUN yarn config set registry https://registry.npm.taobao.org/
RUN yarn && yarn build

FROM nginx:alpine
LABEL MAINTAINER="system"

COPY ./docker_nginx.conf /etc/nginx/conf.d/docker_nginx.conf
COPY --from=0 /redirect_echart/dist /usr/share/nginx/html
RUN cat /etc/nginx/nginx.conf
RUN cat /etc/nginx/conf.d/docker_nginx.conf
RUN ls -al /usr/share/nginx/html