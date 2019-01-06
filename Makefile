make build-api: clean
	CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o bin/main ./api

make docker: build-api
	docker build -t quirk-api .

make test all:
	go test ./...

make clean:
	rm -rf bin/*