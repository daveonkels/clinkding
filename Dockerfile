FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY clinkding /usr/local/bin/clinkding

ENTRYPOINT ["clinkding"]
