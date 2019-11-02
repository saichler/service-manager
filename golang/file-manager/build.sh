rm -rf ./plugin
mkdir ./plugin
go build -buildmode=plugin -o ./plugin/FileService.so Plugin.go
cp ./plugin/FileService.so ../../../plugins/.
