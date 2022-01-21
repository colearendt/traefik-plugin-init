VERSION := "0.0.1"

build:
	docker build -t colearendt/traefik-plugin-init:{{VERSION}} .

install:
	go install traefik-plugin-init

run:
	go run traefik-plugin-init

run-dev:
    go run main.go
