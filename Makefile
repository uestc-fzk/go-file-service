# 变量
DOCKER_USERNAME=uestcfzk
IMAGE_NAME=fileservice
VERSION=v2

swag-init: # swag文档重新生成
	swag init -g ./cmd/main.go -o ./docs
go-mod:
	go mod download
	go mod tidy
go-start: go-mod
	go run ./cmd/main.go

go-build-to-linux:  # 交叉编译
	SET CGO_ENABLED=0 &&\
	SET GOOS=linux&&\
	SET GOARCH=amd64&&\
	go build ./cmd/main.go

# 下面是在Linux环境下进行的
docker-build-image:
	docker build -t $(DOCKER_USERNAME)/$(IMAGE_NAME):$(VERSION) .

docker-run: # 这里挂载一个数据卷rm