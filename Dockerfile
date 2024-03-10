FROM alpine:3.14.0

WORKDIR /app
RUN mkdir -p data/hd css

COPY css/output.css css/
COPY tmp/main ./main
COPY data/symbols.csv data/

# ENTRYPOINT ["/bin/sh","-c","sleep infinity"]
CMD ["./main"]
