FROM dkron/dkron:v2.0.0 as dkron

FROM alpine

COPY --from=dkron /opt/local/dkron /opt/local/dkron
COPY dkron-executor-nats /etc/dkron/plugins/

EXPOSE 8080 8946

ENV SHELL /bin/bash
WORKDIR /opt/local/dkron

ENTRYPOINT ["/opt/local/dkron/dkron"]

CMD ["--help"]
