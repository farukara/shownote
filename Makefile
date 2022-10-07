PHONY: release clean install
dist/dstask: clean
	go build -mod=vendor -o gotask main.go

release:
	./do-release.sh

clean:
	rm -rf dist

install:
	cp dist/dstask /usr/local/bin
	cp dist/dstask-import /usr/local/bin

test:
	go test -v -mod=vendor ./...
	./integrationtest.sh | cat  # cat -- no tty, no confirmations
lint:
	"qa/lint.sh"

update_deps:
	go get
	go mod vendor
	git add -f vendor
