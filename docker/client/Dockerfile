ARG VERSION
ARG IMAGE_PREFIX=tdengine/tdengine
ARG ARTIFACTS=${IMAGE_PREFIX}-artifacts:${VERSION}
FROM ${ARTIFACTS}
LABEL MAINTAINER="Huo Linhe <linhe.huo@gmail.com>"
RUN ls /usr/src/

ARG IMAGE_PREFIX=tdengine/tdengine
FROM ${IMAGE_PREFIX}-runtime
ARG VERSION
ENV VERSION=${VERSION}
COPY --from=0 /usr/src/TDengine-client* /usr/src/
RUN tar xf TDengine-client*.tar.gz && cd TDengine-client-*/ && ./install_client.sh
RUN rm -rf TDengine*
