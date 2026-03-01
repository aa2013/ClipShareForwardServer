@echo off
setlocal enabledelayedexpansion

set "filePath=..\web\web.go"
for /f "tokens=2 delims==" %%a in ('findstr /c:"const version = " "%filePath%"') do (
    set "ver=%%a"
)
:: 去掉双引号
set "ver=%ver:"=%"
:: 去掉空格
set "ver=%ver: =%"
echo Version is: %ver%

:: 编译
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1
set CC=x86_64-w64-mingw32-gcc
set CXX=x86_64-w64-mingw32-g++
go build -ldflags="-s -w" -o ../build/windows/forward_server_windows_amd64_v%ver%.exe ../

endlocal