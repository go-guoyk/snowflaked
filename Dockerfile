FROM golang:1.13 AS builder
ENV GOPROXY https://goproxy.io
WORKDIR /workspace
ADD . .
RUN go build

FROM debian:buster
COPY --from=builder /workspace/snowflaked /usr/local/bin/snowflaked
CMD ["snowflaked"]
