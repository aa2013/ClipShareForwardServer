filePath="../web/web.go"
version=$(awk -F'"' '/const version/ {print $2}' "$filePath")
echo "Version is: $version"

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=1
go build -ldflags="-s -w" -o "../build/linux/forward_server_linux_amd64_v${version}" ../

export GOARCH=arm64
export CC=aarch64-linux-gnu-gcc
export CXX=aarch64-linux-gnu-g++
go build -ldflags="-s -w" -o "../build/linux/forward_server_linux_arm64_v${version}" ../