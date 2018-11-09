SSH_USER ?= test
SSH_HOST ?= ssh-1.mc.lolipop.jp
SSH_PORT ?= 10022
ROTATE ?= 5

build:
	go build

deploy:
	$(eval SUFFIX = $(shell date '+%Y%m%d-%H%M%S'))
	$(eval KEEP = $(shell expr ${ROTATE} + 1))
	env GOOS=linux GOARCH=amd64 go build -o myapp-${SUFFIX} main.go
	scp -P ${SSH_PORT} ./myapp-${SUFFIX} ${SSH_USER}@${SSH_HOST}:/var/app/myapp-${SUFFIX}
	ssh -p ${SSH_PORT} ${SSH_USER}@${SSH_HOST} 'ln -sf /var/app/myapp-${SUFFIX} /var/app/myapp'
	ssh -p ${SSH_PORT} ${SSH_USER}@${SSH_HOST} 'ls -t /var/app/myapp-* | tail -n+${KEEP} | xargs rm'
