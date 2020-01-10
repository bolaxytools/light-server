bd=build
pkg=serverpkg.tar.gz

build:
	mkdir -p $(bd)/bin
	GOOS=linux GOARCH=amd64 go build -o $(bd)/bin/LightServer
	cp config.toml log4go.xml $(bd)/bin/
	tar zcf $(pkg) $(bd)/bin

clean:
	rm -rf $(bd)
	rm -rf $(pkg)

rebuild:
	make clean
	make build