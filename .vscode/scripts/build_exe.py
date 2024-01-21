import json
import os
import winsound
# 環境変数 WORKSPACE_FOLDER の値を取得
workspace_folder = os.environ.get('WORKSPACE_FOLDER')

# 値を表示
print(f"workspace_folder: {workspace_folder}")

# Read app_config.json file
with open(f'{workspace_folder}/cmd/resources/app_config.json', 'r') as file:
    config_data = json.load(file)

# Convert JSON data to dictionary
config_dict = dict(config_data)

app_name = config_dict.get('AppName')
app_version = config_dict.get('AppVersion')

print(f"app_name: {app_name}")
print(f"app_version: {app_version}")

# Build command
# FYNE_THEME=light ライトテーマでビルド
# -o 出力フォルダ
# -trimpath ビルドパスを削除
# -v ビルドログを出力
# -a 全ての依存関係を再ビルド
# -buildmode=exe 実行可能ファイルを生成
# -ldflags "-s -w" バイナリサイズを小さくする
# -H=windowsgui コンソールを表示しない
# -gcflags "all=-N -l" デバッグ情報を削除
build_command = f"go build -o {workspace_folder}/build/{app_name}_{app_version}.exe -trimpath " \
                f"-v -a -buildmode=exe -ldflags \"-s -w -H=windowsgui -X main.Version={app_version}\" {workspace_folder}/cmd/main.go"

print(f"build_command: {build_command}")

os.system(build_command)

# Play beep sound
winsound.PlaySound("SystemAsterisk", winsound.SND_ALIAS)
