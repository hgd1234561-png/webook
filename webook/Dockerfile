FROM ubuntu:20.04

# 把编译后的打包进这个镜像，放到目录/app中
COPY webook /app/webook
WORKDIR /app

ENTRYPOINT ["/app/webook"]