module newC64

go 1.18

require (
	github.com/Djoulzy/Tools/clog v0.0.0-20220609190146-71af779f6ddc
	github.com/Djoulzy/Tools/confload v0.0.0-20220609190146-71af779f6ddc
	github.com/Djoulzy/emutools/mem v0.0.0-20220624083055-f5d43f4b7324
	github.com/Djoulzy/emutools/mos6510 v0.0.0-20220624083055-f5d43f4b7324
	github.com/Djoulzy/emutools/render v0.0.0-20220624083055-f5d43f4b7324
	github.com/mattn/go-tty v0.0.4
)

require (
	github.com/Djoulzy/emutools/charset v0.0.0-20220624083055-f5d43f4b7324 // indirect
	github.com/albenik/bcd v0.0.0-20170831201648-635201416bc7 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-ini/ini v1.66.6 // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/veandco/go-sdl2 v0.4.24 // indirect
	golang.org/x/exp v0.0.0-20220613132600-b0d781184e0d // indirect
	golang.org/x/image v0.0.0-20220617043117-41969df76e82 // indirect
	golang.org/x/sys v0.0.0-20220704084225-05e143d24a9e // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace (
	github.com/Djoulzy/emutools/mem v0.0.0-20220511172747-82eb75258b66 => ../github.com/Djoulzy/emutools/mem
	github.com/Djoulzy/emutools/mos6510 v0.0.0-20220511172747-82eb75258b66 => ../github.com/Djoulzy/emutools/mos6510
	github.com/Djoulzy/emutools/render v0.0.0-20220511172747-82eb75258b66 => ../github.com/Djoulzy/emutools/render
)
