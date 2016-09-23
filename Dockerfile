FROM alpine:3.4
MAINTAINER Oliver Eilhard <oliver@eilhard.net>
ADD docker-test-web /docker-test-web
ENTRYPOINT ["/docker-test-web"]
