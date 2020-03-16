APP_NAME="mockserver"
IMGNAME="ncatelli/${APP_NAME}"
PKG="github.com/PacketFire/${APP_NAME}"

build: | generate
	go build

generate:
	go generate ./...

build-docker: | generate fmt test
	docker build -t ${IMGNAME}:latest .

test: | generate
	go test -race -cover ./...

benchmark:
	go test -benchmem -bench . ./...

fmt:
	test -z $(shell go fmt ./...)

clean-docker:
	@type docker >/dev/null 2>&1 && \
	docker rmi -f ${IMGNAME}:latest || \
	true

clean: clean-docker
	@rm -f ${APP_NAME} || true
	@rm ./pkg/router/generator/plugins.go

lint:
	golint -set_exit_status ./...
