ARG BASE_IMAGE=ubuntu:20.04
FROM ${BASE_IMAGE}
ARG APT_MIRROR=http://mirrors.aliyun.com
ENV DEBIAN_FRONTEND=noninteractive
WORKDIR /usr/src/
RUN sed -Ei "s#http://.*.ubuntu.com#${APT_MIRROR}#" /etc/apt/sources.list \
    && apt-get update \
    && apt-get install -y locales tzdata curl wget \
    && locale-gen en_US.UTF-8 \
    && apt-get clean
ENV LC_ALL=en_US.UTF-8
ENV LANG=en_US.UTF-8
ENV LANGUAGE=en_US.UTF-8
COPY third-party/* /usr/bin/
