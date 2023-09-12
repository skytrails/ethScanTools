PROJECT:=go-admin

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -a -installsuffix "" -o go-admin .

# make build-linux
build-linux-arm64:
	@docker build --platform linux/arm64 -t go-admin:latest .
	@echo "build successful"

# make build-linux
build-linux:
	@docker build --platform linux/amd64 -t go-admin:latest .
	@echo "build successful"

build-sqlite:
	go build -tags sqlite3 -ldflags="-w -s" -a -installsuffix -o go-admin .

# make run
run:
    # delete go-admin-api container
	@if [ $(shell docker ps -aq --filter name=go-admin --filter publish=8000) ]; then docker rm -f go-admin; fi

	# 启动方法一 run go-admin-api container  docker-compose 启动方式
	# 进入到项目根目录 执行 make run 命令
	#@docker-compose up -d

	# 启动方式二 docker run  这里注意-v挂载的宿主机的地址改为部署时的实际决对路径
	@docker run --name=go-admin -p 8000:8000 -v ./config:/go-admin-api/config  -v ./static:/go-admin-api/static -v ./temp:/go-admin-api/temp -d --restart=always go-admin:latest

	@echo "go-admin service is running..."

	# delete Tag=<none> 的镜像
	@docker image prune -f
	@docker ps -a | grep "go-admin"

stop:
    # delete go-admin-api container
	@if [ $(shell docker ps -aq --filter name=go-admin --filter publish=8000) ]; then docker-compose down; fi
	#@if [ $(shell docker ps -aq --filter name=go-admin --filter publish=8000) ]; then docker rm -f go-admin; fi
	#@echo "go-admin stop success"


#.PHONY: test
#test:
#	go test -v ./... -cover

#.PHONY: docker
#docker:
#	docker build . -t go-admin:latest

#save
save:
	docker save -o ./go-admin.tar go-admin:latest

package:

	make build-linux
	make save

# make deploy
deploy:

	#@git checkout master
	#@git pull origin master
	make build-linux
	make run
