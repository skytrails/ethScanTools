PROJECT:=eth-scan

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -a -installsuffix "" -o eth-scan .

# make build-linux
build-linux-arm64:
	@docker build --platform linux/arm64 -t eth-scan:latest .
	@echo "build successful"

# make build-linux
build-linux:
	@docker build --platform linux/amd64 -t eth-scan:latest .
	@echo "build successful"

build-sqlite:
	go build -tags sqlite3 -ldflags="-w -s" -a -installsuffix -o eth-scan .

# make run
run:
    # delete eth-scan-api container
	@if [ $(shell docker ps -aq --filter name=eth-scan --filter publish=8000) ]; then docker rm -f eth-scan; fi

	# 启动方法一 run eth-scan-api container  docker-compose 启动方式
	# 进入到项目根目录 执行 make run 命令
	#@docker-compose up -d

	# 启动方式二 docker run  这里注意-v挂载的宿主机的地址改为部署时的实际决对路径
	@docker run --name=eth-scan -p 8000:8000 --link mysql:mysql -v ./config:/config  -v ./static:/eth-scan-api/static -v ./temp:/eth-scan-api/temp -d --restart=always eth-scan:latest

	@echo "eth-scan service is running..."

	# delete Tag=<none> 的镜像
	@docker image prune -f
	@docker ps -a | grep "eth-scan"

stop:
    # delete eth-scan-api container
	@if [ $(shell docker ps -aq --filter name=eth-scan --filter publish=8000) ]; then docker-compose down; fi
	#@if [ $(shell docker ps -aq --filter name=eth-scan --filter publish=8000) ]; then docker rm -f eth-scan; fi
	#@echo "eth-scan stop success"


#.PHONY: test
#test:
#	go test -v ./... -cover

#.PHONY: docker
#docker:
#	docker build . -t eth-scan:latest

#save
save:
	docker save -o ./eth-scan.tar eth-scan:latest

package:

	make build-linux
	make save

# make deploy
deploy:

	#@git checkout master
	#@git pull origin master
	make build-linux
	make run
