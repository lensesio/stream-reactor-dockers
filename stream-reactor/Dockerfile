ARG STREAM_REACTOR_VERSION
FROM streamreactor/stream-reactor-base:${STREAM_REACTOR_VERSION}

ARG ARCHIVE
ARG URL

ENV ARCHIVE=${ARCHIVE}
ENV URL=${URL}

RUN wget ${URL} && tar -xf ${ARCHIVE} -C /opt/lenses/lib

CMD ["dumb-init", "/opt/lenses/bin/entry-point"]
