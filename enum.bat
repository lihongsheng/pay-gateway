@echo off
setlocal enabledelayedexpansion

REM 统计proto文件数量
set count=0
for /f %%a in ('dir /b "enum\*.proto" 2^>nul ^| find /c /v ""') do set count=%%a

if %count% equ 0 (
    echo not find proto file
    exit /b 0
)

echo find %count% proto files，gen...

REM 批量处理所有proto文件
set success=0
set fail=0

for %%f in (enum\*.proto) do (
    echo.
    echo handle proto file: %%f
    protoc --proto_path=. ^
                --proto_path=./third_party ^
                --go_out=paths=source_relative:. ^
                --go-errors_out=paths=source_relative:. ^
        "%%f"

    if %errorlevel% equ 0 (
        echo success: %%~nf.proto
        set /a success+=1
    ) else (
        echo fail: %%~nf.proto
        set /a fail+=1
    )
)

echo.
echo ======================
echo total:
echo success: %success%
echo fail: %fail%
echo ======================

endlocal
