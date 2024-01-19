@echo off
setlocal enabledelayedexpansion

REM バージョン番号を読み出す

:: JSONファイルのパスを指定
set JSON_FILE=%WORKSPACE_FOLDER%/cmd/resources/app_config.json

:: PowerShellを使用してJSONからデータを抽出
for /f "tokens=1* delims==" %%a in ('powershell -Command "Get-Content %JSON_FILE% | ConvertFrom-Json | ForEach-Object { $$_PSObject.Properties.Name)=$$_PSObject.Properties.Value) }"') do (
    set "%%a=%%b"
)

:: 抽出したデータを表示
echo workspace: %WORKSPACE_FOLDER%
echo AppName: %AppName%
echo AppVersion: %AppVersion%

pause
exit 0

REM ライトモード
set FYNE_THEME=light

REM Windows 64bit
set GOOS=windows
set GOARCH=amd64

REM -o 出力フォルダ
REM -trimpath ビルドパスを削除
REM -v ビルドログを出力
REM -a 全ての依存関係を再ビルド
REM -buildmode=exe 実行可能ファイルを生成
REM -ldflags "-s -w" バイナリサイズを小さくする
REM -gcflags "all=-N -l" デバッグ情報を削除
go build -o %WORKSPACE_FOLDER%/build/%AppName%_%AppVersion%.exe -trimpath -v -a -buildmode=exe -ldflags "-s -w -H=windowsgui -X 'main.Version=%AppVersion%'" %WORKSPACE_FOLDER%/cmd/main.go

msg * "ビルド完了"

endlocal
