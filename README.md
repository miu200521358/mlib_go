# mlib_go

## 新規作成

```
(mtool) C:\MMD\mlib_go\backend>go mod init github.com/miu200521358/mlib_go
go: creating new go.mod: module github.com/miu200521358/mlib_go
```

## Clone Build

- VSCode に go 拡張をいれる `golang.go`
- コマンドプロンプトで下記を実行（× Power Shell）
- https://jmeubank.github.io/tdm-gcc/download/ をインストール
    - `Minimal online installer. ` (多分最初のexeボタン)
    - PATH に TDM-GCC-64 までのパスを追加
- go.mod の最後の3行(replaceのとこ)をコメントアウト
    - 開発用にローカルのを参照しているため
- ライブラリインポート
    - `set GOOS=windows`
    - `set GOARCH=amd64`
    - `set CGO_ENABLED=1`
    - `go mod tidy`
- ビルドできるはず
    - `set CGO_ENABLED=1`
    - `go run main.go`
    - `go build -o ../build/mlib.exe main.go`
- ビルドできたらVSCode再起動で run もできるようになってるはず

## キャッシュクリア系

```
go clean --modcache
go clean -cache
go clean -testcache
```

## フォーマット

```
go fmt ./...
```

## アイコンの組み込み

1. .icoファイルを作成
2. .rcファイルを作成
    - `IDI_ICON1 ICON DISCARDABLE "app.ico"`
3. .resファイルにコンパイル
    - `windres -O coff -o app.res app.rc`
4. ビルドスクリプトで実行


## プロファイル

1. `go run crumb/profile.go`
2. `go tool pprof crumb\profile.go crumb\cpu.pprof`
    - `go tool pprof -flat crumb\profile.go crumb\cpu.pprof`
    - `go tool pprof -cum crumb\profile.go crumb\cpu.pprof`
    - `go tool pprof -http=:8080 cpu.pprof`
3. `(pprof) top`
4. プロファイル: ISAOミク+Addiction


### プロファイルのビジュアライザ

1. `go get github.com/goccy/go-graphviz/cmd/dot`
2. `go install github.com/goccy/go-graphviz`
4. `go tool pprof -http=:8081 mem.pprof`


## SWIG/Bullet

### 前提条件
1. SWIG インストール: https://rinatz.github.io/swigdoc/abstract.html
2. TDM-GCC インストール: https://jmeubank.github.io/tdm-gcc/download/

### SWIG コード生成コマンド

新しいパッケージ構成（`pkg\infra\drivers\mbullet\bt`）用:

```bash
cd pkg\infra\drivers\mbullet\bt

swig_bt.bat
```

### 生成されるファイル
- `bt.go` - Goバインディング
- `bt.cxx` - C++ラッパー
- `bt.h` - ヘッダーファイル

### 注意事項
- `bullet.i` を変更した場合のみ再実行が必要
- SWIG の再生成はリポジトリ管理者が実施
- パスは環境に応じて調整すること

## バージョン反映

```
go list -m -mod=mod -versions github.com/miu200521358/dds
go list -m -mod=mod -versions github.com/miu200521358/win
go list -m -mod=mod -versions github.com/miu200521358/walk
```


### GCプロファイル

1. `set GOGC=1000`
2. `set GODEBUG=gctrace=1`
3. `(mtool) C:\MMD\mlib_go\crumb>go run profile.go`
4. `go tool pprof profile.go cpu.pprof`
5. `go tool pprof -http=:8080 cpu.pprof`

```
export filename=pprof_model_load_20260201_201437_000
printf "top 50\n" | go tool pprof "$filename.pprof" > "${filename}_top_50.txt"
```

### Agent Skills

```
conda create -n mlib python=3.14 -y
conda activate mlib
pip install skillport
skillport init
```

```
レビュワーからの指摘に対する検討および対応計画を立案してください
```

```

```

## License

Source code is licensed under CC-BY-NC-4.0.
Official binaries distributed by the author may be used commercially
under the LICENSE-EXCEPTION.
