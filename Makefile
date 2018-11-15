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
	$(GO) build -o server

server: build
	./server

deploy:
	$(eval SUFFIX = $(shell date '+%Y%m%d-%H%M%S'))
	$(eval KEEP = $(shell expr ${ROTATE} + 1))
	rm -rf ./tmp.dist
	mkdir tmp.dist/
	env GOOS=linux GOARCH=amd64 $(GO) build -o server
	echo DB_NAME=${DB_NAME} >> ./tmp.dist/.env
	echo DB_USER=${DB_USER} >> ./tmp.dist/.env
	echo DB_PASS=${DB_PASS} >> ./tmp.dist/.env
	echo DB_HOST=${DB_HOST} >> ./tmp.dist/.env
	cp -r ${RELEASE_FILES} tmp.dist/
	ssh -p ${SSH_PORT} ${SSH_USER}@${SSH_HOST} 'mkdir -p /var/app/releases'
	scp -r -P ${SSH_PORT} ./tmp.dist ${SSH_USER}@${SSH_HOST}:/var/app/releases/${SUFFIX}
	ssh -p ${SSH_PORT} ${SSH_USER}@${SSH_HOST} 'rm /var/app/current'
	ssh -p ${SSH_PORT} ${SSH_USER}@${SSH_HOST} 'ln -sf /var/app/releases/${SUFFIX} /var/app/current'
	ssh -p ${SSH_PORT} ${SSH_USER}@${SSH_HOST} 'ls -t /var/app/releases/* | tail -n+${KEEP} | xargs --no-run-if-empty rm -rf'
	rm server

initdb:
	scp -P ${SSH_PORT} ./initdb.sql ${SSH_USER}@${SSH_HOST}:/tmp/initdb.sql
	ssh -p ${SSH_PORT} ${SSH_USER}@${SSH_HOST} 'MYSQL_PWD=${DB_PASS} mysql --host=${DB_HOST} --user=${DB_USER} -D ${DB_NAME} < /tmp/initdb.sql'

dev-start:
	echo DB_NAME=${DB_NAME} > ./.env
	echo DB_USER=${DB_USER} >> ./.env
	echo DB_PASS=${DB_PASS} >> ./.env
	echo DB_HOST=127.0.0.1 >> ./.env
	docker run -it --rm --name mc-go-mysql -e MYSQL_ROOT_PASSWORD=${DB_PASS} -e MYSQL_DATABASE=${DB_NAME} -e MYSQL_USER=${DB_USER} -e MYSQL_PASSWORD=${DB_PASS} -v ${PWD}:/docker-entrypoint-initdb.d -d -p 3306:3306 mysql

dev-stop:
	docker stop mc-go-mysql

logs-out:
	ssh -p ${SSH_PORT} ${SSH_USER}@${SSH_HOST} "find /var/log/container/ -type f -name haconiwa.out | xargs ls --full-time | sort -k6,7 | awk '{print \$$9}' | xargs -i cat {}"

logs-err:
	ssh -p ${SSH_PORT} ${SSH_USER}@${SSH_HOST} "find /var/log/container/ -type f -name haconiwa.err | xargs ls --full-time | sort -k6,7 | awk '{print \$$9}' | xargs -i cat {}"
