server:
	go run main.go

swagger:
	swag init --md ./

.PHONY: server swagger
