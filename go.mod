module github.com/miu200521358/mlib_go

go 1.25.5

replace github.com/miu200521358/walk => ../walk

replace github.com/miu200521358/win => ../win

replace github.com/miu200521358/dds => ../dds

require (
	github.com/tiendc/go-deepcopy v1.7.2
	golang.org/x/text v0.33.0
	gonum.org/v1/gonum v0.16.0
)

require golang.org/x/tools v0.40.0 // indirect
