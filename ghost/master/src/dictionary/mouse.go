package dictionary

import (
	"strconv"

	shiori "github.com/Narazaka/shiorigo"
)

type MouseTalks []map[string]Talks

func tutukiTalk() MouseTalks {
	ret := make(MouseTalks, charaCount)

	ret[0] = map[string]Talks{
		"": menu(),
		"Head": {
			"\\0\\\\0側 頭つつき反応",
		},
	}

	ret[1] = map[string]Talks{
		"": {
			"\\1\\\\1側 汎用つつき反応",
		},
	}

	return ret
}

func nadeTalk() MouseTalks {
	ret := make(MouseTalks, charaCount)

	ret[0] = map[string]Talks{
		"": {
			"\\0\\\\0側 汎用なで反応",
		},
		"Head": {
			"\\0\\\\0側 頭なで反応",
		},
	}

	ret[1] = map[string]Talks{
		"": {
			"\\1\\\\1側 汎用なで反応",
		},
	}

	return ret
}

func onMouseDoubleClick(req shiori.Request, vars *Variables) (shiori.Response, error) {
	part := req.Reference(4)
	scope, err := strconv.Atoi(req.Reference(3))
	if err != nil {
		return ResponseNoContent(), err
	}

	if talks, ok := tutukiTalk()[scope][part]; ok {
		return ResponseOneOf(talks), nil
	}

	return ResponseNoContent(), nil
}

func onMouseMove(req shiori.Request, vars *Variables) (shiori.Response, error) {
	part := req.Reference(4)
	scope, err := strconv.Atoi(req.Reference(3))
	if err != nil {
		return ResponseNoContent(), err
	}

	if vars.LastNadePart == part {
		vars.MouseMoveCount[scope]++
	} else {
		vars.MouseMoveCount[scope] = 0
		vars.LastNadePart = part
	}

	if vars.MouseMoveCount[scope] >= mouseMoveLimit {
		vars.MouseMoveCount[scope] = 0
		vars.MoveResetCount[scope] = 0
		return ResponseOneOf(nadeTalk()[scope][part]), nil
	}
	return ResponseNoContent(), nil
}
