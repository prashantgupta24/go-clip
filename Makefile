all: open

build: clean
	mkdir -p -v ./bin/go-clip.app/Contents/Resources
	mkdir -p -v ./bin/go-clip.app/Contents/MacOS
	cp ./appInfo/*.plist ./bin/go-clip.app/Contents/Info.plist
	cp ./appInfo/*.icns ./bin/go-clip.app/Contents/Resources/icon.icns
	go build -o ./bin/go-clip.app/Contents/MacOS/go-clip cmd/main.go
build-win:
	env GOOS=windows GOARCH=amd64 go build -ldflags -H=windowsgui -o ./bin/go-clip-win.exe cmd/main.go
open: build
	open ./bin

clean:
	rm -rf ./bin

start:
	go run cmd/main.go

vet:
	go vet ./...

lint:
	go list ./... | grep -v vendor | grep -v /assets/ |xargs -L1 golint -set_exit_status

test:
	go test -v -failfast ./...
#Race