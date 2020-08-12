env GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(cat .version)" -o ./bin/esdeploy
env GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(cat .version)" -o ./bin/esdeploy.exe
