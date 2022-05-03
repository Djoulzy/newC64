module newC64

go 1.17

require (
	github.com/Djoulzy/Tools/clog v0.0.0-20220429054701-4c221b41ecdf
	github.com/Djoulzy/Tools/confload v0.0.0-20220429054701-4c221b41ecdf
	github.com/Djoulzy/emutools/mem v0.0.0-20220503143932-c59c96d764c0
	github.com/Djoulzy/emutools/mos6510 v0.0.0-20220503143932-c59c96d764c0
	github.com/Djoulzy/emutools/render v0.0.0-20220503143932-c59c96d764c0
	github.com/mattn/go-tty v0.0.4
)

require (
	github.com/Djoulzy/emutools/charset v0.0.0-20220503143932-c59c96d764c0 // indirect
	github.com/albenik/bcd v0.0.0-20170831201648-635201416bc7 // indirect
	github.com/go-ini/ini v1.66.4 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/veandco/go-sdl2 v0.4.20 // indirect
	golang.org/x/image v0.0.0-20220413100746-70e8d0d3baa9 // indirect
	golang.org/x/sys v0.0.0-20220502124256-b6088ccd6cba // indirect
)

replace (
	github.com/Djoulzy/emutools/mem v0.0.0-20220503143932-c59c96d764c0 => ../github.com/Djoulzy/emutools/mem
	github.com/Djoulzy/emutools/mos6510 v0.0.0-20220503143932-c59c96d764c0 => ../github.com/Djoulzy/emutools/mos6510
	github.com/Djoulzy/emutools/render v0.0.0-20220503143932-c59c96d764c0 => ../github.com/Djoulzy/emutools/render
)
