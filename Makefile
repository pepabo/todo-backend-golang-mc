SSH_USER ?= test
SSH_HOST ?= ssh-1.mc.lolipop.jp
SSH_PORT ?= 10022
DB_NAME ?= dbname
DB_USER ?= user
DB_PASS ?= pass
DB_HOST ?= mysql-1.mc.lolipop.lan
ROTATE ?= 5

GO ?= GO111MODULE=on go

RELEASE_FILES = server templates

build:
	$(GO) build

deploy:
	$(eval SUFFIX = $(shell date '+%Y%m%d-%H%M%S'))
	$(eval KEEP = $(shell expr ${ROTATE} + 1))
	rm -rf ./dist
	mkdir dist/
	env GOOS=linux GOARCH=amd64 $(GO) build -o server main.go
	echo DB_NAME=${DB_NAME} >> ./dist/.env.production
	echo DB_USER=${DB_USER} >> ./dist/.env.production
	echo DB_PASS=${DB_PASS} >> ./dist/.env.production
	echo DB_HOST=${DB_HOST} >> ./dist/.env.production
	cp -r ${RELEASE_FILES} dist/
	ssh -p ${SSH_PORT} ${SSH_USER}@${SSH_HOST} 'mkdir -p /var/app/releases'
	scp -r -P ${SSH_PORT} ./dist ${SSH_USER}@${SSH_HOST}:/var/app/releases/${SUFFIX}
	ssh -p ${SSH_PORT} ${SSH_USER}@${SSH_HOST} 'ln -sf /var/app/releases/${SUFFIX} /var/app/current'
	ssh -p ${SSH_PORT} ${SSH_USER}@${SSH_HOST} 'ls -t /var/app/releases/* | tail -n+${KEEP} | xargs --no-run-if-empty rm -rf'
	rm server

initdb:
	scp -P ${SSH_PORT} ./initdb.sql ${SSH_USER}@${SSH_HOST}:/tmp/initdb.sql
	ssh -p ${SSH_PORT} ${SSH_USER}@${SSH_HOST} 'MYSQL_PWD=${DB_PASS} mysql --host=${DB_HOST} --user=${DB_USER} -D ${DB_NAME} < /tmp/initdb.sql'
