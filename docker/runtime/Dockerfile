ARG BASE_IMAGE=ubuntu:20.04
FROM ${BASE_IMAGE}
ARG APT_MIRROR=http://mirrors.aliyun.com
ENV DEBIAN_FRONTEND=noninteractive
WORKDIR /usr/src/
RUN sed -Ei "s#http://.*.ubuntu.com#${APT_MIRROR}#" /etc/apt/sources.list \
    && apt-get update \
    && apt-get install -y locales tzdata curl wget \
    && locale-gen en_US.UTF-8 \
    && wget -c https://github.com/BurntSushi/ripgrep/releases/download/12.1.1/ripgrep_12.1.1_amd64.deb \
    && apt-get install ./ripgrep_12.1.1_amd64.deb \
    && wget -c https://github.com/lotabout/rargs/releases/download/v0.3.0/rargs-v0.3.0-x86_64-unknown-linux-musl.tar.gz \
    && tar xvf rargs-v0.3.0-x86_64-unknown-linux-musl.tar.gz -C /usr/bin \
    && apt-get clean
ENV LC_ALL=en_US.UTF-8
ENV LANG=en_US.UTF-8
ENV LANGUAGE=en_US.UTF-8