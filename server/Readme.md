### 安装说明,以linux为例
> 安装golang环境(略)

> 修改env环境
  ```blade
    go env -w GOPROXY=https://goproxy.cn,direct
  ```
> 修改config.ini
  * redis 密码
  * postgresql 用户名和密码
  * share 分享地址，这个是前端的域名

>  编译
 ```blade
   go mod init
   go mod tidy
   go build -o redirect-admin ./cmd/admin/main.go
   go build -o redirect-endpoint ./cmd/endpoint/main.go  
 ```  
> 运行
 ```blade
   ./redirect-admin
   ./redirect-endpoint
  ```
