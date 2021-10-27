ARG VERSION
ARG IMAGE_PREFIX=tdengine/tdengine
ARG ARTIFACTS=${IMAGE_PREFIX}-artifacts:${VERSION}
FROM ${ARTIFACTS}
LABEL MAINTAINER="Huo Linhe <linhe.huo@gmail.com>"
RUN ls /usr/src/

ARG IMAGE_PREFIX=tdengine/tdengine
FROM ${IMAGE_PREFIX}-runtime

COPY --from=0 /usr/src/TDengine-server* /usr/src/
RUN tar xf TDengine-server*.tar.gz && cd TDengine-server-*/ && ./install.sh -e no
RUN rm -rf TDengine*

WORKDIR /etc/taos
EXPOSE 6030 6031 6032 6033 6034 6035 6036 6037 6038 6039 6040 6041 6042
CMD ["taosd"]
VOLUME [ "/var/lib/taos", "/var/log/taos", "/corefile" ]
ENTRYPOINT [ "/usr/bin/entrypoint.sh" ]
