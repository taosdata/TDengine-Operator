ARG VERSION
ARG IMAGE_PREFIX=tdengine/tdengine
ARG ARTIFACTS=${IMAGE_PREFIX}-artifacts:${VERSION}
FROM ${ARTIFACTS}
LABEL MAINTAINER="Huo Linhe <linhe.huo@gmail.com>"
RUN ls /usr/src/

ARG IMAGE_PREFIX=tdengine/tdengine
FROM ${IMAGE_PREFIX}-runtime

COPY --from=0 /usr/src/TDengine-arbitrator* /usr/src/
RUN tar xf TDengine*-arbitrator*.tar.gz \
  && cd TDengine*-arbitrator-*/ \
  && ./install_arbi.sh \
  && rm -rf TDengine*

EXPOSE 6042/tcp

WORKDIR /var/log/taos/

CMD ["tarbitrator"]

ENTRYPOINT [ "/usr/bin/entrypoint.sh" ]
