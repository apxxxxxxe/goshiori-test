package dictionary

import (
	shiori "github.com/Narazaka/shiorigo"
)

func onKeyPress(req shiori.Request, vars *Variables) (shiori.Response, error) {
	key := req.Reference(0)
	switch key {
	case "r":
		return ResponseOK("\\![reload.shiori]\\0辞書をリロードしました"), nil
	default:
		return ResponseOK("\\0\\s[0]" + key + "キーが押されました"), nil
	}
}
