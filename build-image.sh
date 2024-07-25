#!/bin/sh
# eg.1 : sh build-image.sh
# eg.2, set image: sh build-image.sh zjutjh/jh-forum

# 获取版本信息并确保小写
VERSION=$(git describe --tags --always | cut -f1,2 -d "-" | tr '[:upper:]' '[:lower:]') # eg.: 0.2.5
IMAGE="zjutjh/jh-forum" # 确保镜像名称是小写的

if [ -n "$1" ]; then
  IMAGE=$(echo "$1" | tr '[:upper:]' '[:lower:]') # 确保镜像名称是小写的
fi
if [ -n "$2" ]; then
  VERSION=$(echo "$2" | tr '[:upper:]' '[:lower:]') # 确保版本信息是小写的
fi

# 打印版本信息和镜像名称
echo "Building Docker image with the following details:"
echo "IMAGE: $IMAGE"
echo "VERSION: $VERSION"

# 构建镜像
docker build \
  --build-arg USE_DIST="yes" \
  --tag "$IMAGE:${VERSION}" \
  --tag "$IMAGE:latest" \
  . -f Dockerfile

# 推送到镜像仓库（可选）
#if [ $? -eq 0 ]; then
#  echo "Docker image built successfully."
  # 如果需要推送镜像，请取消以下注释
  # docker push "$IMAGE:${VERSION}"
  # docker push "$IMAGE:latest"
#else
#  echo "Docker image build failed."
#  exit 1
#fi
