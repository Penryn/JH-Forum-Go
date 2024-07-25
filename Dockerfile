# syntax=docker/dockerfile:experimental

# 最终阶段
FROM bitbus/paopao-ce-backend-builder:latest

# 设置环境变量
ARG API_HOST
ARG USE_API_HOST=yes
ARG EMBED_UI=yes
ARG USE_DIST=no
ENV TZ=Asia/Shanghai

# 设置工作目录
WORKDIR /app/jh-forum

# 从构建上下文中复制应用程序和配置文件
COPY release/JH-Forum .
COPY config.yaml.example config.yaml

# 设定数据卷
VOLUME ["/app/data/custom"]

# 开放端口
EXPOSE 8008

# 健康检查
HEALTHCHECK --interval=5s --timeout=3s --retries=3 CMD [ -f /app/jh-forum/JH-Forum ] && ps -ef | grep JH-Forum || exit 1

# 设置容器入口点和默认命令
ENTRYPOINT ["/app/jh-forum/JH-Forum"]
CMD ["serve"]
