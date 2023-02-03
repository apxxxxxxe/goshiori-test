package dictionary

import (
	"encoding/json"
	"os"
	"path/filepath"

	shiori "github.com/Narazaka/shiorigo"
	"github.com/apxxxxxxe/goshiori-test/aozoragetter/balloon"
)

// 定数
const (
	name           = "goshiori"
	version        = "0.1.0"
	charaCount     = 2
	varFile        = "var.json"
	mouseMoveLimit = 50
	moveResetLimit = 10
)

// グローバル変数
// 終了時に保存の必要がない(起動ごとに初期化する)ものは`json:"-"`
type Variables struct {
	DicDir              string          `json:"-"`
	SecondsFromLastTalk int             `json:"-"`
	TalkFrequency       int             `json:"talk_frequency"`
	MouseMoveCount      []int           `json:"-"`
	LastNadePart        string          `json:"-"`
	MoveResetCount      []int           `json:"-"`
	RubyOption          int             `json:"ruby_option"`
	Balloon             balloon.Balloon `json:"-"`
}

type RequestHandler func(shiori.Request, *Variables) (shiori.Response, error)

type Dictionary struct {
	Handlers  map[string]RequestHandler
	Variables *Variables
}

func SaveVariables(vars Variables) error {
	j, err := json.MarshalIndent(&vars, "", "\t")
	if err != nil {
		return err
	}

	exec, err := os.Executable()
	if err != nil {
		return err
	}
	dicDir := filepath.Dir(exec)

	if err := os.WriteFile(filepath.Join(dicDir, varFile), j, 0644); err != nil {
		return err
	}

	return nil
}

func New() (*Dictionary, error) {
	vars, err := loadVariables()
	if err != nil {
		return nil, err
	}
	if vars == nil {
		// 初回起動時のみ初期化
		vars = &Variables{
			TalkFrequency: 60,
			RubyOption:    4,
		}
	}

	// 起動ごとに初期化
	vars.SecondsFromLastTalk = 0
	vars.MouseMoveCount = initIntArray(charaCount, 0)
	vars.MoveResetCount = initIntArray(charaCount, 0)
	vars.LastNadePart = ""

	handlers := map[string]RequestHandler{}

	// shioriの情報を格納
	Info := map[string]string{
		"version":   version,
		"name":      name,
		"craftman":  "hinotsumi",
		"craftmanw": "日野つみ",
	}
	for event, value := range Info {
		handlers[event] = CreateGetHandlerOf(value)
	}

	handlers["OnFirstBoot"] = onFirstBoot
	handlers["OnBoot"] = onBoot
	handlers["OnClose"] = onClose

	handlers["OnMouseMove"] = onMouseMove
	handlers["OnMouseDoubleClick"] = onMouseDoubleClick

	handlers["OnSecondChange"] = onSecondChange
	handlers["OnKeyPress"] = onKeyPress

	handlers["OnUserInput"] = onUserInput

	handlers[EventMenuBookQuery] = onMenuInputBookquery
	handlers[EventInputBookQuery] = onInputBookQuery

	handlers[EventOnGetBalloonInfo] = onGetBalloonInfo

	return &Dictionary{Variables: vars, Handlers: handlers}, nil
}

func loadVariables() (*Variables, error) {
	exec, err := os.Executable()
	if err != nil {
		return nil, err
	}
	dicDir := filepath.Dir(exec)

	varPath := filepath.Join(dicDir, varFile)

	if _, err := os.Stat(varPath); err != nil {
		// varFileが存在しない
		return nil, nil
	}

	b, err := os.ReadFile(varPath)
	if err != nil {
		return nil, err
	}

	var vars Variables
	if err := json.Unmarshal(b, &vars); err != nil {
		return nil, err
	}

	return &vars, nil
}
