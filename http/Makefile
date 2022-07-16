export tag=v1.0
root:
	export ROOT=git@github.com:273539918/hellogo.git

build:
	echo "building httpserver binary"
	mkdir -p bin/amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/amd64 .

release: build
	echo "building httpserver container"
	docker build -t canghong/httpserver:${tag}

push:release
	echo "pushing canghong/httpserver"
	docker push canghong/httpserver:v1.0
