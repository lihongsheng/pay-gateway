@echo off
setlocal enabledelayedexpansion


REM 统计proto文件数量
set count=0
for /f %%a in ('dir /b "proto\*.proto" 2^>nul ^| find /c /v ""') do set count=%%a

if %count% equ 0 (
    echo not find proto file
    exit /b 0
)

echo find %count% proto files, gen...

REM 逐个处理每个proto文件
echo handler proto files
for %%f in (proto\*.proto) do (
    echo.
    echo ==============================================
    echo handling proto file: %%f

    REM 执行goctl rpc生成命令（单个文件）
    goctl rpc protoc "%%f" --go_out=. --go-grpc_out=. --zrpc_out=. -I .\third_party\ -I ./

    REM 获取不带扩展名的文件名
    set "filename=%%~nf"

    REM 删除当前目录下对应的.go文件
    if exist "!filename!.go" (
        echo delete: !filename!.go
        del "!filename!.go" /f /q
    ) else (
        echo not find: !filename!.go
    )

    REM 删除etc目录下对应的.yaml文件
    if exist "etc\!filename!.yaml" (
        echo delete: etc\!filename!.yaml
        del "etc\!filename!.yaml" /f /q
    ) else (
        echo not find: etc\!filename!.yaml
    )
)

echo.
echo ==============================================
echo end:
endlocal
