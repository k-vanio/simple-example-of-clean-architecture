infra:
	docker-compose up -d

install:
	go install github.com/swaggo/swag/cmd/swag@latest

up: initSwag
	cd cmd/ordersystem && go run main.go wire_gen.go

initSwag:
	swag init -d cmd/ordersystem,internal/infra/web,internal/usecases
	