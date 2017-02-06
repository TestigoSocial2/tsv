FROM scratch

MAINTAINER bcessa <ben@pixative.com>

ADD htdocs /var/www/htdocs

ADD htdocs/tsv.db /data/tsv.db

ADD tsv /

EXPOSE 7788

ENTRYPOINT ["/tsv", "server"]
