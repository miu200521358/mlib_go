@echo off
rem bake_leg_ik �? exe にビルドす�?(出�?: こ�?�フォルダの bake_leg_ik.exe)
setlocal
cd /d "%~dp0"
set GO=go
if exist "%USERPROFILE%\sdk\go1.25.6\bin\go.exe" set GO=%USERPROFILE%\sdk\go1.25.6\bin\go.exe
"%GO%" build -o bake_leg_ik.exe .
if errorlevel 1 (
    echo build failed
    exit /b 1
)
echo build ok: %~dp0bake_leg_ik.exe
endlocal

