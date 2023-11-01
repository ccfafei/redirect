#!/usr/bin/env bash
ADMIN_API_URL="http://124.222.224.195:9092"
WEB_SHARE_DOMAIN="http://124.222.224.195:3003"
ENDPOINT_API_URL="http://124.222.224.195"

function cleanimages() {
  docker ps -a | grep "Exited" | awk '{print $1 }'|xargs docker rm >/dev/null 2>/dev/null
  docker ps -a | grep "Created" | awk '{print $1 }'|xargs docker rm >/dev/null 2>/dev/null
  docker images | grep "<none>" | awk '{print $3 }'|xargs docker rmi >/dev/null 2>/dev/null
}

function editConfig(){
  sed -i "s|^share_domain = .*|share_domain = ${WEB_SHARE_DOMAIN}|" ../server/docker_config.ini
  sed -i "s|host_domain = .*|host_domain = '${ENDPOINT_API_URL}';|" ../server/assets/stat.js
  sed -i "s|host_domain = .*|host_domain = '${ENDPOINT_API_URL}';|" ../server/assets/client.js
  sed -i "s|^VITE_BASE_URL = .*|VITE_BASE_URL = ${ADMIN_API_URL}|" ../web/admin/.env.production
  sed -i "s|^VITE_BASE_URL = .*|VITE_BASE_URL = ${ADMIN_API_URL}|" ../web/echart/.env.production
}

chmod -R 777 container-data
editConfig
cleanimages
docker-compose -f docker-compose.yml up -d --build

