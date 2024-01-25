# mlib_go

## 新規作成

```
(mtool) C:\MMD\mlib_go\backend>go mod init github.com/miu200521358/mlib_go
go: creating new go.mod: module github.com/miu200521358/mlib_go
```

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
go get -u github.com/ungerik/go3d
go get -u github.com/go-gl/gl/v4.4-core/gl
go get -u github.com/go-gl/glfw/v3.3/glfw
```

```
go clean --modcache
go clean -cache
go clean -testcache
```

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

## アイコンの組み込み

1. .icoファイルを作成
2. .rcファイルを作成
    - `IDI_ICON1 ICON DISCARDABLE "app.ico"`
3. .resファイルにコンパイル
    - `windres -O coff -o app.res app.rc`
4. ビルドスクリプトで実行
