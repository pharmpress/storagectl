FROM scratch

ADD ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
ADD bin/storage3ctl-linux64-static /usr/local/bin/storage3ctl

ENTRYPOINT ["/usr/local/bin/storage3ctl"]
