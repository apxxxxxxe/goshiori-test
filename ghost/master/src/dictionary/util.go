package dictionary

import (
	"math/rand"
	"time"

	shiori "github.com/Narazaka/shiorigo"
)

type Talks []string

func (values Talks) OneOf() string {
	list := []string{}
	for _, value := range values {
		if value != "" {
			list = append(list, value)
		}
	}
	length := len(list)
	if length <= 0 {
		return ""
	}
	i := rand.Intn(length)
	return list[i]
}

var defaultHeaders = shiori.ResponseHeaders{
	"Charset": "UTF-8",
}

func ResponseNoContent() shiori.Response {
	return shiori.Response{Protocol: shiori.SHIORI, Version: "3.0", Code: 204, Headers: defaultHeaders}
}

func ResponseBadRequest() shiori.Response {
	return shiori.Response{Protocol: shiori.SHIORI, Version: "3.0", Code: 400, Headers: defaultHeaders}
}

func ResponseInternalServerError() shiori.Response {
	return shiori.Response{Protocol: shiori.SHIORI, Version: "3.0", Code: 500, Headers: defaultHeaders}
}

func ResponseOK(value string) shiori.Response {
	res := shiori.Response{Protocol: shiori.SHIORI, Version: "3.0", Code: 200, Headers: defaultHeaders}
	if value != "" {
		res.Headers["Value"] = value
	}
	return res
}

func ResponseOneOf(values Talks) shiori.Response {
	v := values.OneOf()
	if v != "" {
		return ResponseOK(v)
	}
	return ResponseNoContent()
}

func CreateGetHandlerOf(value string) RequestHandler {
	return func(req shiori.Request, vars *Variables) (shiori.Response, error) {
		if req.Method == shiori.GET && value != "" {
			return ResponseOK(value), nil
		}
		return ResponseNoContent(), nil
	}
}

func ResetRNG() {
	rand.Seed(time.Now().UnixNano())
}

func initIntArray(n, i int) []int {
	ret := make([]int, n)
	for j := range ret {
		ret[j] = i
	}
	return ret
}
