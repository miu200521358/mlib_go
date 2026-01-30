# AGENTS

- 設計/命名/エラー/テストは `mlib-guiding-principles` に従う（`mlib_skills\skills\00_project\010_mlib_guiding_principles\references\001_guiding_principles.md`）
- 再構築フェーズは `mlib-rebuild-procedure` に従う（`mlib_skills\skills\00_project\020_mlib_rebuild_procedure\references\001_rebuild_procedure.md`）
- 出力・コメント・回答は日本語で行う
- ファイルは UTF-8 で作成/更新する
- コメント/メッセージログは意図と処理状況が分かる程度に十分記載する
- 公開ドキュメント/ログにドライブレター付き絶対パスは書かず、相対/マスク表記を使う
- build に含まれない test 以外の内部処理ファイルは `internal/` 配下に置く
- 空の Go パッケージは作らず、必要になった時だけ追加する
- コメントは日本語で記述する
- 全関数にドキュメントコメントを必須とする（テスト関数は除く）
- 複雑なロジックには説明コメントを追加する
- `// 指示: miu200521358` は日本語コメント付与のため削除しない
- Go ファイルを新規作成したら必ず `mlib_go_t4/internal/scripts/add_instruction_header.sh` を実行して `// 指示: miu200521358` を付与し、非ASCII化する
- 推測で断言せず、Unknown は残し、確認に必要なファイル/シンボル/手順を列挙すること。
- 実装完了報告前に、コンパイルエラーがない状態（起動できる状態）を必ず確認・達成する
- 実装完了報告前に、文言ヌケモレチェック（`mlib_skills\skills\00_project\060_mlib_i18n_key_checks\scripts\check_i18n_keys.py`）を必ず実施・ヌケモレを防ぐ
- 特に指定がない場合、ルートディレクトリは /mnt/c/Codex/mlib とする
