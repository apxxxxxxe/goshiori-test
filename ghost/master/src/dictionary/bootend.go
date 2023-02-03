package dictionary

import (
	shiori "github.com/Narazaka/shiorigo"
)

// 起動時共通処理
var commonScript = "\\0\\b[2]\\![get,property," + EventOnGetBalloonInfo + "," +
	"currentghost.balloon.name," +
	"currentghost.balloon.scope(0).lines," +
	"currentghost.balloon.scope(0).validwidth," +
	"currentghost.balloon.scope(0).validheight.initial," +
	"currentghost.balloon.scope(0).validheight," +
	"currentghost.balloon.scope(0).char_width," +
	"currentghost.balloon.scope(0).basepos.y]"

func onFirstBoot(req shiori.Request, vars *Variables) (shiori.Response, error) {
	talks := Talks{
		"\\1\\s[10]\\0\\s[1]初回起動。\\e",
	}

	for i := range talks {
		talks[i] = commonScript + talks[i]
	}

	return ResponseOneOf(talks), nil
}

func onBoot(req shiori.Request, vars *Variables) (shiori.Response, error) {
	talks := Talks{
		"\\1\\s[10]\\0\\s[1]これは Go で栞を作るサンプルだから、\\w2過度な期待はしないよーに。\\e",
	}

	for i := range talks {
		talks[i] = commonScript + talks[i]
	}

	return ResponseOneOf(talks), nil
}

func onClose(req shiori.Request, vars *Variables) (shiori.Response, error) {
	talks := Talks{
		"\\1\\s[10]\\0\\s[5]じゃ、\\w2えんいー！\\_w[500]\\e",
	}
	return ResponseOneOf(talks), nil
}
