package dictionary

import (
	shiori "github.com/Narazaka/shiorigo"
)

const EventInputBookQuery = "OnInputBookQuery"
const EventMenuBookQuery = "OnMenuInputBookQuery"

func menu() Talks {
	talks := Talks{
		"\\0\\b[2]メニュー台詞1。\\n" +
			"\\n" +
			"\\![*]\\q[本を読んでもらう," + EventMenuBookQuery + "]",
	}

	return talks
}

func onMenuInputBookquery(req shiori.Request, vars *Variables) (shiori.Response, error) {
	return ResponseOK("\\![open,inputbox," + EventInputBookQuery + ",timeout=-1,]\\0本のタイトルや作者を教えて。\\n(検索ワードをスペース区切りで指定してください)"), nil
}
