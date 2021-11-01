module github.com/Sora233/DDBOT

go 1.16

require (
	github.com/Jeffail/gabs/v2 v2.6.1
	github.com/Logiase/MiraiGo-Template v0.0.0-20210228150851-29092d4d5486
	github.com/Mrs4s/MiraiGo v0.0.0-20211005153339-2cdb7407f907
	github.com/PuerkitoBio/goquery v1.8.0
	github.com/Sora233/sliceutil v0.0.0-20210120043858-459badd8d882
	github.com/alecthomas/kong v0.2.17
	github.com/davecgh/go-spew v1.1.1
	github.com/ericpauley/go-quantize v0.0.0-20200331213906-ae555eb2afa4
	github.com/go-playground/validator/v10 v10.9.0 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/guonaihong/gout v0.2.6
	github.com/hashicorp/golang-lru v0.5.4
	github.com/json-iterator/go v1.1.12
	github.com/jupp0r/go-priority-queue v0.0.0-20160601094913-ab1073853bde
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/mattn/go-isatty v0.0.13 // indirect
	github.com/modern-go/gls v0.0.0-20190610040709-84558782a674
	github.com/nfnt/resize v0.0.0-20180221191011-83c6a9932646
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/sirupsen/logrus v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/buntdb v1.2.7
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210903071746-97244b99971b // indirect
	golang.org/x/text v0.3.7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace (
	github.com/Logiase/MiraiGo-Template => github.com/Sora233/MiraiGo-Template v0.0.0-20210815095536-29373b391593
	github.com/tidwall/gjson => github.com/tidwall/gjson v1.9.3
	github.com/willf/bitset v1.2.0 => github.com/bits-and-blooms/bitset v1.2.0
)
