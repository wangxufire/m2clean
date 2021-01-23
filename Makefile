BINARY_NAME=m2clean
LDFLAGS := -s -w

darwin:
	env GOOS=darwin GOARCH=amd64 go build -v -trimpath -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)_darwin_amd64

linux:
	env GOOS=linux GOARCH=amd64 go build -v -trimpath -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)_linux_amd64

win32:
	env GOOS=windows GOARCH=amd64 go build -v -trimpath -ldflags "$(LDFLAGS)" -o dist/$(BINARY_NAME)_win32_amd64.exe

clean:
	rm -rf ./dist