# AGENTS ガイド

## このドキュメントについて
- mlib_go は MMD 系 3D ツールのための基盤ライブラリで、ネイティブ API や外部アセットとの橋渡しが多いリポジトリです。
- 本書はエージェントや新規参加者が作業前に押さえておきたいリポジトリの特徴・癖・レイヤ構成をまとめています。
- クリーンアーキテクチャを志向していますが完全分離ではないため、依存方向と責務を意識しながら開発してください。

## リポジトリの性格
- 対応 OS は基本的に Windows（`//go:build windows` が存在）で、CGO を有効にしたビルドが前提です。
- Bullet 物理（`pkg/infrastructure/bt`）や OpenGL（`pkg/infrastructure/mgl`）などのネイティブ連携が多く、開発環境のセットアップが複雑になりがちです。
- モデル・モーション（PMX/VMD）や物理ステート、描画周りの処理をドメインとして扱い、インフラ層で外部ライブラリへ橋渡ししています。
- プロファイル・サンプルなど開発補助用ディレクトリ（`crumb/`, `test_resources/`, `archive/`）が同居しているため、用途を見極めて利用してください。

## レイヤとディレクトリの目安
- `cmd/` : エントリポイントと UI/ユースケース起動コード。用途ごとに `app/`, `ui/`, `usecase/` へ分離。
- `pkg/config/` : ロギング・多言語化・プロセス制御など横断的な設定群。
- `pkg/domain/` : MMD ドメインの中心。`pmx`, `vmd`, `physics`, `rendering`, `delta`, `mmath` 等がモデル・値オブジェクト・差分ロジックを提供し、テストもここに集中しています。
- `pkg/usecase/` : ユースケースの調整層。現状は `deform` など限定的ですが、ドメインと UI/インフラを結びます。
- `pkg/interface/` : プレゼンテーション層。`viewer` や `controller` が UI イベントとユースケースの結線を担当します。
- `pkg/infrastructure/` : 外部技術との接続。`mbt`（Bullet ラッパーと風力・デバッグ表示）、`mgl`（OpenGL レンダリング）、`mfile`（ファイル IO）、`render`（描画サブシステム）などが存在します。
- `build/` : `go build` で生成した成果物の配置先、`distribution/` は配布物テンプレート群です。

## 主なサブシステムの覚書
- 物理: `pkg/infrastructure/mbt` が Bullet を包み、`pkg/domain/physics` の設定値や `pkg/domain/delta` の差分適用と協調。風シミュレーション（`wind.go`）やデバッグ描画（`debug_view.go`）が同梱されています。
- レンダリング: `pkg/infrastructure/mgl` がシェーダー・バッファの構築を扱い、`pkg/domain/rendering` のステートと同期します。MSAA 設定や floor レンダラーなどが分割管理されています。
- データストア: `pkg/infrastructure/repository` に PMX/VMD 等のロード実装が集まっており、`pkg/domain/model` との変換を担当します。
- サポートユーティリティ: `pkg/infrastructure/mstring`, `pkg/infrastructure/miter` など標準ライブラリ補助がまとまっています。多用すると依存方向が崩れやすいので注意。

## ビルド & 実行の基本
- 前提ツール: Go 1.23 系、TDM-GCC などの MinGW、SWIG（Bullet ラッパー再生成時）。必要な PATH と `CGO_ENABLED=1`, `GOOS=windows`, `GOARCH=amd64` を環境にセットしてください。
- ローカルモジュール: `go.mod` の `replace ../walk`, `../win`, `../dds` はローカル開発用。フォークでリモートを使う場合はコメントアウトを検討します。
- コマンド例: `go fmt ./...`, `go clean --modcache`, `go run main.go`, `go build -o build/mlib.exe cmd/app/main.go` など README の手順に準拠してください。
- プロファイル: `crumb/profile.go` を入口に `go tool pprof` を回す運用。CPU/GPU の観点を整理するため、`pkg/profile` ではなく `crumb/` を使います。

## テストと検証
- ユニットテストは `pkg/domain` 配下に点在（`*_test.go`）。LLRB や幾何計算などドメインロジックを中心に整備されています。
- 物理や描画はネイティブ連携が強く自動テストが難しいため、サンプルプロジェクトとデバッグ描画での目視確認が主手段です。
- 新規機能追加時は差分型（`pkg/domain/delta`）が関わるかを確認し、状態遷移のテストを優先してください。

## よくある落とし穴
- CGO まわり: 環境変数を忘れると Bullet ラッパーがリンクエラーになります。`bt` サブディレクトリを更新した場合は SWIG の再実行も必要です。
- 依存方向: Clean Architecture の層を崩すような import を避け、`domain -> usecase -> interface` の依存関係を逆転させないようにしてください。
- アセット配置: `test_resources/` は大きめの PMX/VMD を含むため、CI で扱う際はパスや容量に注意します。
- Windows 固有 API (`pkg/infrastructure/win` など) に依存するコードはビルドタグでガードされています。新規追加時もタグを忘れずに。

## 今後の整備方針メモ
- ユースケース層の拡充とインターフェース層の依存解消が課題です。新規処理を追加する場合はまずユースケースを定義し、インフラ実装は `factory` 等で注入する形を推奨します。
- 既存の `cmd/` 構成は機能別に肥大化しやすいので、共通初期化コードは `pkg/config` や `pkg/infrastructure/factory` へ寄せる運用で整理してください。
- ドキュメント化されていないビルドスクリプトや手順を発見した場合は README か本ファイルへ追記し、エージェントが再現可能な状態を維持しましょう。
