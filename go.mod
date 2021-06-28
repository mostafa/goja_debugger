module github.com/mostafa/goja_debugger

go 1.16

require (
	github.com/c-bata/go-prompt v0.2.6
	github.com/dop251/goja v0.0.0-20210614154742-14a1ffa82844
	github.com/dop251/goja_nodejs v0.0.0-20210225215109-d91c329300e7
	github.com/mattn/go-colorable v0.1.8 // indirect
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	golang.org/x/sys v0.0.0-20210616094352-59db8d763f22 // indirect
)

replace github.com/dop251/goja => github.com/mstoykov/goja v0.0.0-20210628100343-fe3ffd18f06a
