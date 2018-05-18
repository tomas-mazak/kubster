IMAGE=kubster

build: image

build/kubster: kubster.go
	mkdir -p build
	docker run --rm -v ${PWD}/:/go/src/kubster golang /bin/sh -c "cd /go/src/kubster && go get -v -d . && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/kubster ."

image: build/kubster Dockerfile.imageonly
	docker build --no-cache -f Dockerfile.imageonly -t $(IMAGE) .
