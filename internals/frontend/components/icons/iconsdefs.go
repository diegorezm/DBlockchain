package icons

import (
	"github.com/a-h/templ"
)

var Bolt = Icon("bolt")
var Clock = Icon("clock")
var Wallet = Icon("wallet")
var Pickaxe = Icon("pickaxe")
var Copy = Icon("copy")
var CopyCheck = Icon("copy-check")
var ChevronLeft = Icon("chevron-left")
var RefreshCW = Icon("refresh-cw")
var Save = Icon("save")
var HandCoins = Icon("hand-coins")

func GetIconFromString(name string) func(...Props) templ.Component {
	_, ok := internalSvgData[name]
	if ok {
		return Icon(name)
	}
	return Icon("bolt")
}
