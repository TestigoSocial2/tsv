default: build

# Current version
version = 0.3.1

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
	yarn install
	yarn build
	env GOOS=linux GOARCH=amd64 go build -v -o tsv -ldflags "\
	-X github.com/transparenciamx/tsv/cmd.infoBuild=`git log --pretty=format:'%H' -n1` \
	-X github.com/transparenciamx/tsv/cmd.infoVersion=$(version)"
	docker build -t tm/tsv:$(version) .
	docker save -o tsv.tar tm/tsv
	gzip tsv.tar
	rm tsv

# Run as a docker container
run-docker:
	docker run -it -p 7788:7788 --rm \
	-v /var/run/tsv-pre-register:/data/pre-register \
	tm/tsv server -s /data -c /var/www/htdocs

## Generate the list of valid CA certificates
ca-roots:
	@docker run -dit --rm --name caroots ubuntu:16.04
	@docker exec --privileged caroots sh -c "apt update"
	@docker exec --privileged caroots sh -c "apt install -y ca-certificates"
	@docker exec --privileged caroots sh -c "cat /etc/ssl/certs/* > /ca-certificates.crt"
	@docker cp caroots:/ca-certificates.crt ca-certificates.crt
	@docker stop caroots