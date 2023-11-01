##
## Build 
##
FROM golang:1.18-alpine AS redirect_admin_builder
WORKDIR /app
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env \
    && go mod tidy \
    && go build -o redirect-admin ./cmd/admin/main.go

##
## Deploy
##
FROM alpine:latest
WORKDIR /app

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && apk del tzdata \
    && mkdir logs \
    && chown -R 777 logs

COPY --from=redirect_admin_builder /app/redirect-admin .
COPY --from=redirect_admin_builder /app/docker_config.ini .

EXPOSE 9092
ENTRYPOINT ./redirect-admin -c docker_config.ini


