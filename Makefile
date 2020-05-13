build-mac:
	GOOS=darwin GOARCH=amd64 go build

build-v1-mac:
	GOOS=darwin GOARCH=amd64 go build -o compare-tables.v1

