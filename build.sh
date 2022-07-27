GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o diskio-windows-amd64.exe .
GOOS=linux   GOARCH=amd64 go build -ldflags "-s -w" -o diskio-linux-amd64.exe   .
GOOS=darwin  GOARCH=amd64 go build -ldflags "-s -w" -o diskio-darwin-adm64.exe  .
GOOS=darwin  GOARCH=arm64 go build -ldflags "-s -w" -o diskio-darwin-arm64.exe  .