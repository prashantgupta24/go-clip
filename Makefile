all: open

build: clean
	mkdir -p -v ./bin/go-clip.app/Contents/Resources
	mkdir -p -v ./bin/go-clip.app/Contents/MacOS
	cp ./appInfo/*.plist ./bin/go-clip.app/Contents/Info.plist
	cp ./appInfo/*.icns ./bin/go-clip.app/Contents/Resources/icon.icns
	go build -o ./bin/go-clip.app/Contents/MacOS/go-clip systray/main.go

open: build
	open ./bin

clean:
	rm -rf ./bin

start:
	go run systray/main.go

vet:
	go vet $(shell glide nv)

lint:
	go list ./... | grep -v vendor | grep -v /assets/ |xargs -L1 golint -set_exit_status
