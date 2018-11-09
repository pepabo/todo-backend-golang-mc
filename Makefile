SSH_USER ?= test
SSH_HOST ?= ssh-1.mc.lolipop.jp
SSH_PORT ?= 10022

build:
	env GOOS=linux GOARCH=amd64 go build -o myapp main.go

deploy: build
	scp -P ${SSH_PORT} ./myapp ${SSH_USER}@${SSH_HOST}:/var/app
