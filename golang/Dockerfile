FROM alpine:3.14.0

WORKDIR /app
RUN mkdir -p data/hd

COPY main .
COPY data/symbols.csv data/


CMD ["./main"]