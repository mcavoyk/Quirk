make build-api:
	CGO_ENABLED=0 GOOS=linux go build -installsuffix cgo -o bin/main ./api

make docker: build-api
	docker build -t quirk-api .

