package dictionary

import (
	"fmt"
	"regexp"

	shiori "github.com/Narazaka/shiorigo"
	"github.com/apxxxxxxe/goshiori-test/aozoragetter/book"
)

var repSpaces = regexp.MustCompile(`[ 　]`)

func onUserInput(req shiori.Request, vars *Variables) (shiori.Response, error) {
	id := req.Reference(0)
	text := req.Reference(1)

	talks := Talks{
		"\\![raise," + id + "," + text + "]",
	}

	return ResponseOneOf(talks), nil
}

func onInputBookQuery(req shiori.Request, vars *Variables) (shiori.Response, error) {
	text := req.Reference(0)

	talks, err := searchAndGetBook(text, vars)
	if err != nil {
		return ResponseInternalServerError(), err
	}

	return ResponseOneOf(talks), nil
}

func searchAndGetBook(queries string, vars *Variables) (Talks, error) {
	candidates, src, err := book.GetBook(repSpaces.Split(queries, -1))
	if err != nil {
		return nil, err
	}
	switch len(candidates) {
	case 0:
		return Talks{
			"……そういう本はないみたいだ。",
		}, nil
	case 1:
		vars.Balloon.Reset()
		vars.Balloon.StoreVirticalWriting(src)
		return Talks{
			"……うん。取ってくるよ。",
		}, nil
	default:
		books := "\\_q\\n\\n"
		for _, c := range candidates {
			books += fmt.Sprintf("\\q[%s「%s」（%s）,%s,%s %s %s]\\n", c.Title, c.Author, c.Mojidukai, EventInputBookQuery, c.ID, c.Title, c.Author)
		}
		books += "\\_q\\_l[0,0]これだけの候補があるみたいだよ。"
		return Talks{
			books,
		}, nil
	}
}
