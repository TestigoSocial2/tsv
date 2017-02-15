default: build

# Current version
version = 0.2.0

# Build for the default architecture in use
build:
	go build -v -ldflags "\
	-X github.com/transparenciamx/tsv/cmd.infoBuild=`git log --pretty=format:'%H' -n1` \
	-X github.com/transparenciamx/tsv/cmd.infoVersion=$(version)"

# Redownload all dependencies
clean:
	rm -rf vendor glide.lock
	glide install

# Run as a native CLI
run:
	./tsv server

# Build server for Linux and pack it with the contents
docker:
	env GOOS=linux GOARCH=amd64 go build -v -o tsv -ldflags "\
	-X github.com/transparenciamx/tsv/cmd.infoBuild=`git log --pretty=format:'%H' -n1` \
	-X github.com/transparenciamx/tsv/cmd.infoVersion=$(version)"
	docker build -t tm/tsv .
	docker save -o tsv.tar tm/tsv
	gzip tsv.tar
	rm tsv

# Landing page build
landing:
	mv htdocs/index.html htdocs/backup.up
	mv htdocs/landing.html htdocs/index.html
	make docker
	mv htdocs/index.html htdocs/landing.html
	mv htdocs/backup.up htdocs/index.html

# Run as a docker container
run-docker:
	docker run -it -p 7788:7788 --rm \
	-v /var/run/tsv-pre-register:/data/pre-register \
	tm/tsv server -s /data -c /var/www/htdocs

# Perform a basic deploy on a local server (run as root)
deploy:
	docker cp tsv:/data/tsv.db /etc/tsv/data.db
	systemctl stop tsv
	gzip -d tsv.tar.gz
	load -i tsv.tar
	rm tsv.tar
	systemctl start tsv
