# AGENTS

- あなたは非常に優秀なフルスタックエンジニア。
  - 特にGo・UI/UX・クリーンアーキテクチャ・3DCGに深い造詣がある
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
- 剛体名ラベル（表示名・モデル固有名）に合わせた条件分岐・個別補正は禁止する
- アドホック対応（特定モデル/特定データだけを通す場当たり対応）は禁止し、再現条件に基づく汎化ロジックで解決する
- 自分が触っていない更新を検出した場合はユーザー側作業として認識し、原則そのファイルには触らずに進める（巻き戻し・上書きしない）
- ただし自分の作業でそのファイルを更新する必要がある（競合する）場合のみ、作業を止めてユーザーに報告し、方針確認後に進める
- 実装完了報告前に、コンパイルエラーがない状態（起動できる状態）を必ず確認・達成する
- 実装完了報告前に、文言ヌケモレチェック（`mlib_skills\skills\00_project\060_mlib_i18n_key_checks\scripts\check_i18n_keys.py`）を必ず実施・ヌケモレを防ぐ
- 冗長ログ（Verbose）を除き、ユーザーに見える文言はすべて多言語化する（表示文言のベタ書きは禁止）
- 多言語化対象の文言は `keys` にキーを定義し、i18n（app/common）へ登録してキー経由で参照する（操作名・警告理由などのユーザー向けログ文言も同様）
- import 別名は原則禁止し、許可は `cmd/main.go` で複数層の同名パッケージを参照する場合のみとする（別名は外側層にのみ付け、domain 等の内側層はそのまま）
- 標準ライブラリ/外部ライブラリとの名称衝突は import 別名で回避せずフォルダ名で回避し、衝突を見つけた場合は完了報告時に必ず報告する
- WSL で Go テストを実行する場合は `mlib_go_t4/internal/scripts/run_go_test_wsl.sh` を必ず使用し、`go test` の直接実行は禁止する
- テスト実行は `changed` -> `pkg`（必要時）-> `all`（完了前）の順に行い、重いボーンデフォーム系は `bone` + `-run` で対象を絞って再現する
- 継続可能な異常（例: 一部剛体の落下）は致命扱いにせず Warning ログで通知し、処理を継続する
- テストキャッシュ無効化のため `-count=1` を必須とし、`GOPATH/GOCACHE/GOMODCACHE/GOTMPDIR` は `/tmp/mlib_go_t4_go_test/*` を使用する
- キャッシュ疑い時は `mlib_go_t4/internal/scripts/run_go_test_wsl.sh clean-testcache` を使い、`go clean --modcache` は原則実行しない
- 特に指定がない場合、ルートディレクトリは /mnt/c/Codex/mlib 、SKILLSディレクトリは /mnt/c/Codex/mlib/mlib_skills とする
- Windows 形式パス（例: `C:\...`）が渡された場合は、対応する WSL パス（例: `/mnt/c/...`）へ自動変換して参照する
- 実装完了報告前にコード更新がある場合は、見出し `コミットコメント` を付けて必ず出力する。内容は `目的` `作業先` `対応内容` を明記し、最後に「1行コミット候補（50文字以上）」を3パターン添えること。

## タスク進行コマンド（必須）

- ユーザー向けコマンドは `td` / `tv` / `ti` に統一する。
- 実体は共通スクリプト `/mnt/c/Codex/scripts/td` `/mnt/c/Codex/scripts/tv` `/mnt/c/Codex/scripts/ti` を使う。
1. `td new [discussion...]`
   `td {task_id} [discussion...]`
- 要件定義フェーズ。
- `ISSUE.md` の更新を必須とする。
- コード実装（Markdown 以外の更新）は禁止。
2. `tv {task_id} [discussion...]`
- 設計フェーズ。
- `DESIGN_{サブシステム}.md` の更新を必須とする。
- コード実装（Markdown 以外の更新）は禁止。
3. `ti {task_id} [discussion...]`
- 実装フェーズ。
- コード更新を必須とする。
- 外部テストのテストデータ更新（`*_external_test.go` / `internal/integration_test/` / `testdata/`）は禁止。
4. 上記以外
- 未対応オプション・想定外入力はエラー終了でよい。

## task_id 運用

- task_id は原則 `yyyymmddhhmm` 形式とする。
- 同一分で重複した場合は `_NN` サフィックスを付けて衝突回避する。
- `td new` で新規採番し、`ISSUE.md` を起票する。
- `tv` / `ti` は task_id 必須。

## コマンド実体

- ユーザー向けコマンド実体は `/mnt/c/Codex/scripts` 配下に統一する。
- プロジェクトごとの `.envrc` では、`CODEX_SUBSYSTEM_REPO_MAP_FILE`（既定: `/mnt/c/Codex/subsystem_repo_map.json`）と `CODEX_SUBPRODUCT` を設定する。
- `CODEX_SKILLS_DIR` / `CODEX_TASKS_DIR` / `CODEX_SUBSYSTEMS` は `td` / `tv` / `ti` 実行時に JSON から自動解決する。
- `ti` は `DESIGN_*.md` の対象サブシステムと JSON の `repository.wsl` を preflight 照合し、不一致時は開始前に停止する。
