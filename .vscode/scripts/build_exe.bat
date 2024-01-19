@echo off
setlocal enabledelayedexpansion

REM �o�[�W�����ԍ���ǂݏo��

:: JSON�t�@�C���̃p�X���w��
set JSON_FILE=%WORKSPACE_FOLDER%/cmd/resources/app_config.json

:: PowerShell���g�p����JSON����f�[�^�𒊏o
for /f "tokens=1* delims==" %%a in ('powershell -Command "Get-Content %JSON_FILE% | ConvertFrom-Json | ForEach-Object { $$_PSObject.Properties.Name)=$$_PSObject.Properties.Value) }"') do (
    set "%%a=%%b"
)

:: ���o�����f�[�^��\��
echo workspace: %WORKSPACE_FOLDER%
echo AppName: %AppName%
echo AppVersion: %AppVersion%

pause
exit 0

REM ���C�g���[�h
set FYNE_THEME=light

REM Windows 64bit
set GOOS=windows
set GOARCH=amd64

REM -o �o�̓t�H���_
REM -trimpath �r���h�p�X���폜
REM -v �r���h���O���o��
REM -a �S�Ă̈ˑ��֌W���ăr���h
REM -buildmode=exe ���s�\�t�@�C���𐶐�
REM -ldflags "-s -w" �o�C�i���T�C�Y������������
REM -gcflags "all=-N -l" �f�o�b�O�����폜
go build -o %WORKSPACE_FOLDER%/build/%AppName%_%AppVersion%.exe -trimpath -v -a -buildmode=exe -ldflags "-s -w -H=windowsgui -X 'main.Version=%AppVersion%'" %WORKSPACE_FOLDER%/cmd/main.go

msg * "�r���h����"

endlocal
