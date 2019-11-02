rm -rf ./plugin
mkdir ./plugin
env GOOS=freebsd GOARCH=amd64 go build -buildmode=plugin -o ./plugin/FileService.so Plugin.go
cp ./plugin/FileService.so ../../../plugins/.
scp ./plugin/FileService.so root@192.168.86.169:/root/.
