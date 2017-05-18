FROM scratch

MAINTAINER bcessa <ben@pixative.com>

ADD ca-certificates.crt /etc/ssl/certs/ca-certificates

ADD htdocs /var/www/htdocs

ADD tsv /

EXPOSE 7788 7789

ENTRYPOINT ["/tsv"]
