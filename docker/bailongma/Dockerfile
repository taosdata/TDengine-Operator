ARG VERSION
ARG IMAGE_PREFIX=tdengine/tdengine
FROM ${IMAGE_PREFIX}-client:${VERSION}
LABEL MAINTAINER="Huo Linhe <linhe.huo@gmail.com>"
ARG BAILONGMA_VERSION
ENV BAILONGMA_VERSION=${BAILONGMA_VERSION}
RUN curl -o /usr/bin/bailongma https://github.com/taosdata/bailongma-rs/releases/download/v${BAILONGMA_VERSION}/bailongma-amd64 \
  && chmod +x /usr/bin/bailongma
COPY bin/* /usr/bin/
ENTRYPOINT [ "/usr/bin/entrypoint.sh" ]
CMD ["bailongma"]
