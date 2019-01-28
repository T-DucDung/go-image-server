VERSION=`git symbolic-ref -q --short HEAD || git describe --tags --exact-match`
CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w -X main.buildVersion=$VERSION -X main.buildTime=`date -Is`"
upx -qqq go-image-server
