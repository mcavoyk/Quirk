build-api: clean
	GOOS=linux go build -o bin/main ./api

build-api-windows: clean
	GOOS=windows go build -o bin/main.exe ./api

docker: build-api
	docker build -t quirk-api .

# End to End testing
test-e2e:
	go test ./api/tests -v

# Unit and integration tests
test:
	go test `go list ./api/... | grep -v tests` -v

test-all:
	go test ./api/... -v

clean:
	rm -rf bin/*