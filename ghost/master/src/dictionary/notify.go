package dictionary

import (
	"strconv"

	shiori "github.com/Narazaka/shiorigo"
	"github.com/apxxxxxxe/gohst/aozoragetter/balloon"
)

const EventOnGetBalloonInfo = "OnGetBalloonInfo"

const (
	fontSize = 14
	fontName = "UD Digi Kyokasho N-R,ＭＳ 明朝"
	margin   = 4
)

func onGetBalloonInfo(req shiori.Request, vars *Variables) (shiori.Response, error) {
	width, err := strconv.Atoi(req.Reference(2))
	if err != nil {
		return ResponseInternalServerError(), err
	}

	height, err := strconv.Atoi(req.Reference(3))
	if err != nil {
		return ResponseInternalServerError(), err
	}

	baseYPos, err := strconv.Atoi(req.Reference(6))
	if err != nil {
		return ResponseInternalServerError(), err
	}

	vars.Balloon = balloon.Balloon{
		Width:     width,
		Height:    height,
		FontSize:  fontSize,
		FontName:  fontName,
		Margin:    margin,
		BaseYPos:  baseYPos,
		XPos:      0,
		YPos:      0,
		PageIndex: 0,
		Pages:     nil,
	}

	return ResponseNoContent(), nil
}
