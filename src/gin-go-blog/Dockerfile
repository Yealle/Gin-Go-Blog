# # 正常镜像
# FROM golang:latest

# ENV GOPROXY https://goproxy.cn,direct

# Scratch小镜像
FROM scratch 

WORKDIR $GOPATH/src/gin-go-blog
COPY . $GOPATH/src/gin-go-blog

# RUN go build .

EXPOSE 8000
# ENTRYPOINT 指令是指定容器启动程序及参数
# 注意是模块名
ENTRYPOINT [ "./gin-blog" ]
