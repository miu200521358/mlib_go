import os
import json


def go_files_to_json(pkg_full_path, output_json_path="mlib.json"):
    go_files_dict = {}

    # 指定されたディレクトリを再帰的に探索
    for root, dirs, files in os.walk(pkg_full_path):
        for file in files:
            # ディレクトリ名 "bt"以下はスキップ。ただし"mbt"以下は処理対象
            if "bt" in root and "mbt" not in root:
                continue
            # 拡張子が.goかつ_testを含まないファイルのみ処理
            if file.endswith(".go") and "_test" not in file:
                # ファイルのフルパスを取得
                full_path = os.path.join(root, file)
                # 相対パスをキーに設定
                relative_path = os.path.relpath(full_path, pkg_full_path)
                # ファイル内容を読み込み
                with open(full_path, 'r', encoding='utf-8') as f:
                    file_lines = f.readlines()
                # 辞書に追加
                go_files_dict[relative_path] = "\n".join(file_lines)

    # 辞書をJSON形式で書き出し
    with open(output_json_path, 'w', encoding='utf-8') as json_file:
        json.dump(go_files_dict, json_file, ensure_ascii=False, indent=4)

    print(f"JSONファイルが {output_json_path} に保存されました。")


# 使用例
pkg_full_path = "C:/MMD/mlib_go/pkg"
go_files_to_json(pkg_full_path)

