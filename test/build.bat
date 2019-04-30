cd ../../../../
set GOPATH=%cd%
set GOARCH=amd64
set GOOS=windows
cd src/demo/melody_test/test
go build -v -ldflags="-s -w"