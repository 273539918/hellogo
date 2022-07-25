# http服务构建过程
cd /golang/examples/gohttpserver
## 构建镜像
make release
## 构建并上传镜像
make push
## 运行容器
docker run -p 8080:80 -d canghong/httpserver:v1.0
## 访问服务
curl localhost:8080/env/version
curl localhost:8080/header
curl localhsot:8080/healthz