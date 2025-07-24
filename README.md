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


## bullet

1. swig インストール
https://rinatz.github.io/swigdoc/abstract.html

2. 変換コード作成

```
(mtool) C:\MMD\mlib_go\pkg\bt>swig -c++ -go -cgo -I"C:\MMD\mlib_go\pkg\bt\bullet\src" -I"C:\development\TDM-GCC-64\lib\gcc\x86_64-w64-mingw32\10.3.0\include\c++\x86_64-w64-mingw32" -I"C:\development\TDM-GCC-64\x86_64-w64-mingw32\include" -I"C:\Program Files\Microsoft Visual Studio\2022\Community\VC\Tools\MSVC\14.38.33130\include" -cpperraswarn -o "C:\MMD\mlib_go\pkg\bt\bt.cxx" "C:\MMD\mlib_go\pkg\bt\bullet.i"
```

## バージョン反映

```
go list -m -mod=mod -versions github.com/miu200521358/dds
go list -m -mod=mod -versions github.com/miu200521358/win
go list -m -mod=mod -versions github.com/miu200521358/walk
```

---------


## fyne (没)

```
go get fyne.io/fyne/v2@latest
go install fyne.io/fyne/v2/cmd/fyne@latest
```

```
C:\MMD\mlib_go>fyne bundle resources\MPLUS1-Regular.ttf > pkg\front\core\bundle.go
go get fyne.io/fyne/v2/internal/svg@v2.4.3
go get fyne.io/fyne/v2/storage/repository@v2.4.3
```

```
go get fyne.io/fyne/v2/internal/driver/glfw@v2.4.3
go get fyne.io/fyne/v2/app@v2.4.3
go get fyne.io/fyne/v2/widget@v2.4.3
go get fyne.io/fyne/v2/internal/painter@v2.4.3
```

```
fyne bundle icon.png > icon.go
```

```
go get fyne.io/fyne/v2
go get github.com/ungerik/go3d
go get github.com/fyne-io/glfw-js
go get fyne.io/fyne/v2/layout
```

### walk

```
go get github.com/akavel/rsrc
cd %GOPATH%\pkg\mod\github.com\akavel\rsrc@v0.10.2
go build
```

```
rsrc -manifest main.manifest -o rsrc.syso
```

```
go get -u golang.org/x/image
```


### GCプロファイル

1. `set GOGC=1000`
2. `set GODEBUG=gctrace=1`
3. `(mtool) C:\MMD\mlib_go\crumb>go run profile.go`
4. `go tool pprof profile.go cpu.pprof`
5. `go tool pprof -http=:8080 cpu.pprof`
