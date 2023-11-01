##
## Build 
##
FROM golang:1.18-alpine AS redirect_endpoint_builder
WORKDIR /app
COPY . .

RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go env -w CGO_ENABLED=0 \
    && go env \
    && go mod tidy \
    && go build -o redirect-endpoint ./cmd/endpoint/main.go

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
    && chown -R 777 logs \
    && mkdir assets

COPY --from=redirect_endpoint_builder /app/redirect-endpoint .
COPY --from=redirect_endpoint_builder /app/docker_config.ini .
COPY --from=redirect_endpoint_builder /app/assets  ./assets

EXPOSE 9091
ENTRYPOINT ./redirect-endpoint -c docker_config.ini

