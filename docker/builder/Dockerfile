ARG BASE_IMAGE=ubuntu:20.04
FROM ${BASE_IMAGE}
LABEL MAINTAINER="Huo Linhe <linhe.huo@gmail.com>"

ARG APT_MIRROR=http://mirrors.aliyun.com
ENV DEBIAN_FRONTEND=noninteractive
RUN sed -Ei "s#http://.*.ubuntu.com#${APT_MIRROR}#" /etc/apt/sources.list \
    && apt-get update \
    && apt-get install -y build-essential cmake wget git\
    && apt-get clean

ARG VERSION
WORKDIR /tmp
COPY script/build.sh /tmp/
RUN /tmp/build.sh

FROM alpine:3
LABEL MAINTAINER="Huo Linhe <linhe.huo@gmail.com>"

ARG VERSION
ENV VERSION=${VERSION:-develop}
COPY --from=0 /tmp/TDengine-*/release/TDengine-* /usr/src/