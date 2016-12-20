FROM scratch

MAINTAINER bcessa <ben@datos.mx>

ADD tmp/htdocs /var/www/htdocs

ADD tmp/tsv.db /data/tsv.db

ADD tsv /

EXPOSE 7788

ENTRYPOINT ["/tsv", "server"]
