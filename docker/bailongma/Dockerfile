ARG VERSION
ARG IMAGE_PREFIX=tdengine/tdengine
FROM ${IMAGE_PREFIX}-client:${VERSION}
LABEL MAINTAINER="Huo Linhe <linhe.huo@gmail.com>"
ARG BAILONGMA_VERSION
ENV BAILONGMA_VERSION=${BAILONGMA_VERSION}
COPY bin/* /usr/bin/
COPY bailongma-${BAILONGMA_VERSION} /usr/bin/bailongma
ENTRYPOINT [ "/usr/bin/entrypoint.sh" ]
CMD ["bailongma"]
