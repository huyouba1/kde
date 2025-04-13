PKG = "github.com/huyouba1/kde"
BUILD_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_COMMIT := ${shell git rev-parse  --short HEAD}
BUILD_TIME := ${shell date '+%Y-%m-%d %H:%M:%S'}
BUILD_GO_VERSION := $(shell go version | grep -o  'go[0-9].[0-9].*')
MAIN_FILE := "cmd/server/main.go"


gen: ## Init Service
	@protoc -I=.  --go_out=. --go_opt=module=${PKG} --go-triple_out=. --go-triple_opt=module=${PKG} api/*/pb/*.proto
	@go fmt ./...
	@protoc-go-inject-tag -input=api/*/*.pb.go


tidy:
	go mod tidy

run:
	@go run cmd/server/main.go

linux:
	@GOOS=linux GOARCH=amd64 go build -a -o build/${OUTPUT_NAME} -ldflags "-s -w" -ldflags "-X '${VERSION_PATH}.GIT_BRANCH=${BUILD_BRANCH}' -X '${VERSION_PATH}.GIT_COMMIT=${BUILD_COMMIT}' -X '${VERSION_PATH}.BUILD_TIME=${BUILD_TIME}' -X '${VERSION_PATH}.GO_VERSION=${BUILD_GO_VERSION}'" ${MAIN_FILE}

