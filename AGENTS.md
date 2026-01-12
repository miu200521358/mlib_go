# AGENTS

- 設計/命名/エラー/テストは `mlib-guiding-principles` に従う（`mlib_skills\skills\00_project\010_mlib_guiding_principles\references\001_guiding_principles.md`）
- 再構築フェーズは `mlib-rebuild-procedure` に従う（`mlib_skills\skills\00_project\020_mlib_rebuild_procedure\references\001_rebuild_procedure.md`）
- 出力・コメント・回答は日本語で行う
- ファイルは UTF-8 で作成/更新する
- コメント/メッセージログは意図と処理状況が分かる程度に十分記載する
- 公開ドキュメント/ログにドライブレター付き絶対パスは書かず、相対/マスク表記を使う
- build に含まれない test 以外の内部処理ファイルは `internal/` 配下に置く
- 空の Go パッケージは作らず、必要になった時だけ追加する
- ローカル/CI のテスト入口は `internal/scripts/test.ps1`
- コメントは日本語で記述する
- 全関数にドキュメントコメントを必須とする（テスト関数は除く）
- 複雑なロジックには説明コメントを追加する
