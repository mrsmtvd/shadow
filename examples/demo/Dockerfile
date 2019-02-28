FROM kihamo/docker-golang-base

EXPOSE 8001
EXPOSE 8080

COPY demo /
ENTRYPOINT ["/demo", "-dashboard.host=0.0.0.0"]