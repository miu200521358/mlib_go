# mlib_go

## 新規作成

```
(mtool) C:\MMD\mlib_go\backend>go mod init github.com/miu200521358/mlib_go
go: creating new go.mod: module github.com/miu200521358/mlib_go
```

## fyne インストール

```
go get fyne.io/fyne/v2@latest
go install fyne.io/fyne/v2/cmd/fyne@latest
```

### 日本語設定

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

