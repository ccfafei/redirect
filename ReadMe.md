## 系统说明及部署

### 1. 项目介绍
> Redirect 是一个基于gin和vue做的跳转管理和统计系统。
  * 支持批量管理域名，支持规则缓存和同步。
  * 支持IP黑名单，在黑名单的IP跳转到默认。
  * 支持轮询、权重均衡设置,有高并发锁。
  * 分享连接查看统计UV,PV,IP;采用HypeLogLog比较节省内存。
  * 支持Redis集群，可以多节点部署跳转程序。
  * 另外还有类似站长统计功能

> 软件环境:
```
- node版本 > v16.8.3
- golang版本 >= v1.16
- postgres版本 >= v14.7
- redis版本 >= v6.2.6
- vue版本   >= v3.0.5
- IDE推荐：Goland

```
> 基本结构说明:

``` blade
- docker 主要部署文件
  |-- docker-compose-admin.yml 部署docker配置
  |-- container-data 容器持久化存放目录，这个要修改可写权限
  |-- build.sh 启动脚本,注意修改变量
- server 主要包含后台和endpoint节点
  |-- admin.Dockerfile     后台dockerfile文件
  |-- endpoint.Dockerfile  节点dockerfile文件
  |-- docker_config.ini    docker部署用到的配置
  |-- config.ini           默认配置，二进制方便调试
  |-- cmd                  admin和endpoint入口main
  |-- core                 核心组件，初始化/中间件/定时器
  |-- router               路由，url路由规则
  |-- controller           控制器，处理请求，分发到服务
  |-- service              服务层，基本业务逻辑
  |-- model                模型，定义表结构及其它结构体
  |-- storage              数据存储，主有redis和postgres数据处理
  |-- utils                工具，一些函数和包的封装
  |-- assets               js统计代码,类似站长统计
  |-- logs                 日志存储目录
  |-- sys_url_redirect.sql  数据库sql
- web 主要是后台有html及展示图表
  |-- admin 后台html
  |   ├── Dockerfile 部署文件
  |   |-- docker_nginx.conf 部署配置
  |   └── .env.production 编译用到配置
  |-- echart 统计html
  |     └── Dockerfile 部署文件
  |          ├── docker_nginx.conf 部署配置
  |          └── .env.production 编译用到配置

 ```
### 2. docker 部署
 > 安装docker,已安装请忽略   
   ```blade
    curl -sSL https://get.docker.com/ | sh
    sudo systemctl start docker #启动docker
    sudo systemctl enable docker #加入开机自启动
   ```
 
 > 安装docker composer，已安装请忽略  
   ```blade
     sudo curl -SL https://github.com/docker/compose/releases/download/v2.17.2/docker-compose-linux-x86_64 -o /usr/local/bin/docker-compose
     sudo chmod +x /usr/local/bin/docker-compose
   ```
 
 > 执行部署    
   ```blade    
      #进入目录      
       cd docker
	  
      #新版docker需要用官方账号登录
	  docker login 用户名
	  输入密码
	 

      # 修改build.sh中变量
      vi build.sh

      # 修改配置信息docker-compose.yml,建议默认,如果修改了对应Dockerfile用到配置也要修改
       vi docker-compose.yml  

      #执行启动脚本
      ./build.sh

      #释放
      docker-compose  -f docker-compose.yml down
   ```
 ### 访问
  ``` blade
   网址: http://localhost:3002/#/login
   用户名: admin
   密码: 123456(默认)
  ```

### 3. linux并发及性能调优
  ``` blade

    设置ulimit,可以根据命令设置
    ulimit -n 262144

    用户级可修改以下配置:
    vi /etc/security/limits.conf

    *   hard    stack   262144
    *   soft     nproc   65536

    设置sysctl
    vi /etc/sysctl.conf

    fs.file-max = 1000000
    net.ipv4.ip_local_port_range = 10000     65000
    net.ipv4.tcp_syncookies = 1
    net.ipv4.tcp_tw_reuse=1 #让TIME_WAIT状态可以重用，这样即使TIME_WAIT占满了所有端口，也不会拒绝新的请求造成障碍 默认是0
    net.ipv4.tcp_tw_recycle=1 #让TIME_WAIT尽快回收 默认0，部分系统不支持这个选项，暂时未弄清楚为何
    net.ipv4.tcp_fin_timeout=30 #时间修改有考究

  ```