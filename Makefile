SSH_USER ?= test
SSH_HOST ?= ssh-1.mc.lolipop.jp
SSH_PORT ?= 10022
DB_NAME ?= dbname
DB_USER ?= user
DB_PASS ?= pass
DB_HOST ?= mysql-1.mc.lolipop.lan
ROTATE ?= 5

GO ?= GO111MODULE=on go

build:
	$(GO) build

deploy:
	$(eval SUFFIX = $(shell date '+%Y%m%d-%H%M%S'))
	$(eval KEEP = $(shell expr ${ROTATE} + 1))
	env GOOS=linux GOARCH=amd64 $(GO) build -o myapp-${SUFFIX} main.go
	scp -P ${SSH_PORT} ./myapp-${SUFFIX} ${SSH_USER}@${SSH_HOST}:/var/app/myapp-${SUFFIX}
	ssh -p ${SSH_PORT} ${SSH_USER}@${SSH_HOST} 'ln -sf /var/app/myapp-${SUFFIX} /var/app/myapp'
	ssh -p ${SSH_PORT} ${SSH_USER}@${SSH_HOST} 'ls -t /var/app/myapp-* | tail -n+${KEEP} | xargs --no-run-if-empty rm'
	rm myapp-${SUFFIX}

initdb:
	scp -P ${SSH_PORT} ./initdb.sql ${SSH_USER}@${SSH_HOST}:/tmp/initdb.sql
	ssh -p ${SSH_PORT} ${SSH_USER}@${SSH_HOST} 'MYSQL_PWD=${DB_PASS} mysql --host=${DB_HOST} --user=${DB_USER} -D ${DB_NAME} < /tmp/initdb.sql'
