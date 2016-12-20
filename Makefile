default: build

# Build for the default architecture in use
build:
	go build -v

# Redownload all dependencies
clean:
	rm -rf vendor glide.lock
	glide install

# Run as a native CLI
run:
	./tsv server -s tmp/ -c tmp/htdocs

# Build server for Linux and pack it with the contents
docker:
	env GOOS=linux GOARCH=amd64 go build -v -o tsv
	docker build -t tm/tsv .
	rm tsv

# Run as a docker container
run-docker:
	docker run -it -p 7788:7788 --rm tm/tsv server -s /data -c /var/www/htdocs
