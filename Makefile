build-api: clean
	CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o bin/main ./api

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