module github.com/miu200521358/mlib_go

go 1.22.4

require (
	github.com/ftrvxmtrx/tga v0.0.0-20150524081124-bd8e8d5be13a
	github.com/go-gl/gl v0.0.0-20231021071112-07e5d0ea2e71
	github.com/go-gl/glfw/v3.3/glfw v0.0.0-20240118000515-a250818d05e3
	github.com/go-gl/mathgl v1.1.0
	github.com/miu200521358/dds v0.0.0
	github.com/miu200521358/walk v0.0.3
	github.com/pkg/profile v1.7.0
	golang.org/x/image v0.15.0
	golang.org/x/text v0.14.0
)

require (
	github.com/felixge/fgprof v0.9.3 // indirect
	github.com/google/pprof v0.0.0-20211214055906-6f57359322fd // indirect
	github.com/miu200521358/win v0.0.1 // indirect
	golang.org/x/exp v0.0.0-20231110203233-9a3e6036ecaa // indirect
	golang.org/x/tools v0.15.0 // indirect
)

require (
	github.com/jinzhu/copier v0.4.0
	github.com/nicksnyder/go-i18n/v2 v2.4.0
	github.com/petar/GoLLRB v0.0.0-20210522233825-ae3b015fd3e9
	golang.org/x/sys v0.14.0
	gonum.org/v1/gonum v0.15.0
	gopkg.in/Knetic/govaluate.v3 v3.0.0 // indirect
)

replace github.com/miu200521358/walk => ../walk

replace github.com/miu200521358/win => ../win

replace github.com/miu200521358/dds => ../dds
