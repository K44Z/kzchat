build:
	go build -o kzchat ./cmd/kzchat/server

install: build
	sudo mv kzchat /usr/local/bin/
