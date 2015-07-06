FROM kihamo/scratch-ca-certs

EXPOSE 8001
EXPOSE 8080

COPY shadow-full /
ENTRYPOINT ["/shadow-full", "-host=0.0.0.0"]