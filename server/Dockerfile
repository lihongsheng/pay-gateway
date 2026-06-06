FROM crpi-enwh9j989b0fsg8b.cn-guangzhou.personal.cr.aliyuncs.com/alexxqli/base:golang_base_1.25.0-v0.0.4 as builder

WORKDIR /backend
# 只有 go.mod 或 go.sum 变化时，这层才会失效
COPY go.mod go.sum ./

# 安装 git
RUN apk add --no-cache git ca-certificates tzdata build-base

# 配置 git 使用用户名密码访问私有仓库
RUN git config --global url."https://lixiangdaqian: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@github.com/".insteadOf "https://github.com/"

RUN go env -w GO111MODULE=on \
    && go env -w CGO_ENABLED=0 \
    && go env -w GOPRIVATE=github.com/singer-stack-lab \
    && go env

# 只有依赖变化时才重新下载，否则使用缓存
RUN go mod download && go mod verify
# 源码变化只影响这层及后续的构建
COPY . .
# 如果只是源码变化，依赖不变，这层会快很多
RUN go build -v -o server .

#RUN go env -w GO111MODULE=on \
#    && go env -w GOPROXY=https://goproxy.cn,direct \
#    && go env -w CGO_ENABLED=0 \
#    && go env \
#    && go mod tidy \
#    && go build -o server .

FROM swr.cn-north-4.myhuaweicloud.com/ddn-k8s/docker.io/library/alpine:latest

LABEL MAINTAINER="SliverHorn@sliver_horn@qq.com"
# 设置时区
ENV TZ=Asia/Shanghai
RUN apk update && apk add --no-cache tzdata openntpd \
    && ln -sf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone


WORKDIR /backend

# 构建参数
ARG ENV=test
COPY --from=0 /backend ./
COPY --from=0 /backend/resource ./resource/
COPY --from=0 /backend/config.${ENV}.yaml ./config.yaml

#COPY --from=0 /go/src/github.com/lihongsheng/pay-gateway/server ./
#COPY --from=0 /go/src/github.com/lihongsheng/pay-gateway/resource ./resource/
#COPY --from=0 /go/src/github.com/lihongsheng/pay-gateway/config.docker.yaml ./

# 挂载目录：如果使用了sqlite数据库，容器命令示例：docker run -d -v /宿主机路径/gva.db:/go/src/github.com/lihongsheng/pay-gateway/gva.db -p 8888:8888 --name gva-server-v1 gva-server:1.0
# VOLUME ["/go/src/github.com/lihongsheng/pay-gateway"]

EXPOSE 8888
ENTRYPOINT ./server -c config.yaml
