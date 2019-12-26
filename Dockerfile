FROM dkron/dkron:v2.0.2 as dkron

FROM alpine

RUN set -x \
	&& buildDeps='bash ca-certificates openssl tzdata' \
	&& apk add --update $buildDeps \
	&& rm -rf /var/cache/apk/* \
	&& mkdir -p /opt/local/dkron

COPY --from=dkron /opt/local/dkron /opt/local/dkron
COPY dkron-executor-nats /etc/dkron/plugins/

EXPOSE 8080 8946

ENV SHELL /bin/bash
WORKDIR /opt/local/dkron

ENTRYPOINT ["/opt/local/dkron/dkron"]

CMD ["--help"]
