FROM golang:latest as builder

# Fetch source from github
RUN go get github.com/tbinhluong/reverseproxy/...
WORKDIR /go/src/github.com/tbinhluong/reverseproxy/
RUN CGO_ENABLED=0 GOOS=linux make build

# Multi-stage build docker image
FROM alpine:latest

LABEL maintainer="Binh Luong <tbinhluong@gmail.com>"

RUN mkdir -p /reverseproxy && \
    chown -R nobody:nogroup /reverseproxy

COPY --from=builder /go/src/github.com/tbinhluong/reverseproxy/dist/reverseproxy /bin/reverseproxy
COPY --from=builder /go/src/github.com/tbinhluong/reverseproxy/config/config.yml /reverseproxy/config.yml

USER nobody
EXPOSE 8080
WORKDIR /reverseproxy
ENTRYPOINT [ "/bin/reverseproxy" ]
CMD [ "--config.file=config.yml" ]
