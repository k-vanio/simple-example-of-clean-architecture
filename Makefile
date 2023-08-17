infra:
	docker-compose up -d

install:
	go install github.com/swaggo/swag/cmd/swag@latest
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go install github.com/ktr0731/evans@latest
	export PATH="$PATH:$(go env GOPATH)/bin"
	printf '// +build tools\npackage tools\nimport (_ "github.com/99designs/gqlgen"\n _ "github.com/99designs/gqlgen/graphql/introspection")' | gofmt > tools.go

genQl:
	go run github.com/99designs/gqlgen generate

initSwag:
	swag init -d cmd/ordersystem,internal/infra/web,internal/usecases

proto:
	protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles/order.proto

up: initSwag
	cd cmd/ordersystem && go run main.go wire_gen.go