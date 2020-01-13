FROM golang:1.13 AS builder
WORKDIR /workspace
ADD . .
RUN go build

FROM debian:buster
ADD docker-entrypoint.sh /
ENTRYPOINT ["/docker-entrypoint.sh"]
COPY --from=builder /workspace/snowflaked /usr/local/bin/snowflaked
CMD ["snowflaked"]
