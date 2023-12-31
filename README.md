# Instruções para Subir Containers

- Este repositório contém um conjunto de definições de containers Docker para criar e gerenciar nossa infraestrutura de desenvolvimento. O processo de configuração é automatizado usando o comando `make infra`.
- Para subir aplicações `make up`.
- Veja documentação [docs](http://localhost:8000/docs/index.html)
## Requisitos

Certifique-se de ter os seguintes requisitos instalados em sua máquina antes de prosseguir:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Golang 1.20](https://go.dev/doc/install)

## Como Usar

1. **Clone o Repositório**

   ```sh
   git clone https://github.com/k-vanio/simple-example-of-clean-architecture.git
   cd simple-example-of-clean-architecture

2. **URL API**

```sh
   http://localhost:8000/api/orders GET
   http://localhost:8000/api/orders POST