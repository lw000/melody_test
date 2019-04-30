cd ../../../../
set GOPATH=%cd%
cd src/demo/melody_test/service
set GOARCH=amd64
set GOOS=windows
go build -v -ldflags="-s -w"