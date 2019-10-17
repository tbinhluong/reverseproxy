FROM golang:latest

# Fetch source from github
RUN go get github.com/tbinhluong/reverseproxy/...
WORKDIR /go/src/github.com/tbinhluong/reverseproxy/
RUN GOOS=linux make build

# Multi-stage build docker image
FROM alpine:latest

LABEL maitainer="Binh Luong <tbinhluong@gmail.com>"

RUN mkdir -p /reverseproxy/config && \
    chown -R nobody:nogroup /reverseproxy

COPY --from=0 /go/src/github.com/tbinhluong/reverseproxy/dist/reverseproxy /reverseproxy/.
COPY --from=0 /go/src/github.com/tbinhluong/reverseproxy/config/config.yml /reverseproxy/config/config.yml

USER nobody
EXPOSE 8080
WORKDIR /reverseproxy
ENTRYPOINT [ "./reverseproxy" ]
CMD [ "--config.file=./config/config.yml" ]