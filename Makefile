server:
	go run main.go

swagger:
	swag init --md ./

test:
	ginkgo -v --cover --coverprofile=coverage.out ./...

coverage:
	go tool cover -func coverage.out

.PHONY: server swagger test coverage
