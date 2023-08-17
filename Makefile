infra:
	docker-compose up -d

install:
	go install github.com/swaggo/swag/cmd/swag@latest
	printf '// +build tools\npackage tools\nimport (_ "github.com/99designs/gqlgen"\n _ "github.com/99designs/gqlgen/graphql/introspection")' | gofmt > tools.go

genQl:
	go run github.com/99designs/gqlgen generate

up: initSwag
	cd cmd/ordersystem && go run main.go wire_gen.go

initSwag:
	swag init -d cmd/ordersystem,internal/infra/web,internal/usecases
	