FROM alpine:3.14.0

WORKDIR /app
RUN mkdir -p data/hd

COPY main ./main
COPY data/symbols.csv data/

# ENTRYPOINT ["/bin/sh","-c","sleep infinity"]
CMD ["./main"]