package icons

import (
	"github.com/a-h/templ"
)

var Bolt = Icon("bolt")
var Clock = Icon("clock")
var Wallet = Icon("wallet")

func GetIconFromString(name string) func(...Props) templ.Component {
	_, ok := internalSvgData[name]
	if ok {
		return Icon(name)
	}
	return Icon("bolt")
}
