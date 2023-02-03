package dictionary

import (
	shiori "github.com/Narazaka/shiorigo"
)

func onSecondChange(req shiori.Request, vars *Variables) (shiori.Response, error) {
	var err error

	// 一定時間後なでカウントをリセットする
	for i := 0; i < charaCount; i++ {
		if vars.MouseMoveCount[i] > 0 {
			vars.MoveResetCount[i]++
		}
		if vars.MoveResetCount[i] >= moveResetLimit {
			vars.MouseMoveCount[i] = 0
			vars.MoveResetCount[i] = 0
		}
	}

	// 一定時間後ランダムトーク
	vars.SecondsFromLastTalk++
	if vars.SecondsFromLastTalk >= vars.TalkFrequency {
		vars.SecondsFromLastTalk = 0
		return ResponseOK(randomTalk()), err
	}
	return ResponseNoContent(), err
}

func randomTalk() string {
	talks := Talks{
		"\\0\\s[0]ランダムトーク1",
		"\\0\\s[1]ランダムトーク2",
		"\\0\\s[2]ランダムトーク3",
	}
	return talks.OneOf()
}
