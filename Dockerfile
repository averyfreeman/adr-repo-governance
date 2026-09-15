FROM alpine:3.22

COPY dist/adr /usr/local/bin/adr

RUN chmod 0755 /usr/local/bin/adr

ENTRYPOINT ["/usr/local/bin/adr"]
