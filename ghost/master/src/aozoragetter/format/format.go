package format

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/apxxxxxxe/goshiori-test/aozoragetter/setting"
	"golang.org/x/net/html"
)

const (
	styleBold = iota
	styleItalic
)

var replastCRLF = regexp.MustCompile(`\n+$`)

func style(text string, style int) string {
	s := ""
	switch style {
	case styleBold:
		s = "bold"
	case styleItalic:
		s = "italic"
	default:
		return text
	}
	return "\\f[" + s + ",1]" + text + "\\f[" + s + ",default]"
}

func expandRuby(ruby *html.Node) string {
	var (
		result string
		rb     string
		rt     string
	)

	for c := ruby.FirstChild; c != nil; c = c.NextSibling {
		switch c.Data {
		case "rb":
			// 漢字
			rb = processNode(c)
		case "rt":
			// よみがな
			rt = processNode(c)
		}
	}

	// どうせこれしか使わないだろうから固定
	option := 4

	switch option {
	case 1:
		// 漢字のみ
		result = rb
	case 2:
		// よみがなのみ
		result = rt
	case 3:
		// 漢字(よみがな)
		result = rb + "(" + rt + ")"
	case 4:
		// 漢字(よみがな) 読み上げはよみがなのみ
		result = "\\_q" + rb + "(\\_q" + rt + "\\_q)\\_q"
	case 5:
		// 漢字(よみがな) 読み上げは漢字のみ
		result = rb + "\\_q(" + rt + ")\\_q"
	default:
		result = rb + "(" + rt + ")"
	}

	return result
}

// http://kumihan.aozora.gr.jp/slabid-19.htmに記載されている注記を処理する
func ParseHtmlText(htmlText string) (string, error) {
	root, err := html.Parse(strings.NewReader(htmlText))
	if err != nil {
		return "", err
	}
	result := processNode(findMainText(root))

	return result, nil
}

func processNode(root *html.Node) string {
	result := ""
	for c := root.FirstChild; c != nil; c = c.NextSibling {
		isProcessed := false
		if c.Type == html.ElementNode {
			switch c.Data {
			case "div":
				isProcessed = true
				if class := extractClass(c); strings.HasPrefix(class, "jisage_") {
					if j, err := strconv.Atoi(strings.TrimPrefix(class, "jisage_")); err == nil {
						child := strings.ReplaceAll(processNode(c), "\n", "\n"+strings.Repeat("　", j))
						result += strings.Repeat("　", j) + strings.TrimRight(child, "　")
					}
				}
			case "span":
				cls := extractClass(c)
				switch cls {
				case "":
					// 未対応
					isProcessed = true
					result += processNode(c)
				case "yokogumi":
					// 未対応
					isProcessed = true
					result += processNode(c)
				case "keigakomi":
					isProcessed = true
					result += "\\n\\n" + processNode(c) + "\\n\\n"
				case "caption":
					// 未対応
					isProcessed = true
					result += "(" + processNode(c) + ")"
				case "warichu":
					isProcessed = true
					result += "〔※" + processNode(c) + "〕"
				case "notes":
					// 注釈は握りつぶす
					child := processNode(c)
					if strings.Contains(child, "改ページ") {
						result += setting.RepageCode
					}
					isProcessed = true
				case "futoji":
					isProcessed = true
					result += style(processNode(c), styleBold)
				case "shatai":
					isProcessed = true
					result += style(processNode(c), styleItalic)
				default:
					if strings.HasPrefix(cls, "sho") || strings.HasPrefix(cls, "dai") {
						isProcessed = true
						result += style(processNode(c), styleBold)
					}
				}
			case "sub":
				switch extractClass(c) {
				case "subscript":
					// 未対応
					isProcessed = true
					result += "(" + processNode(c) + ")"
				case "kaeriten":
					// 未対応
					isProcessed = true
					result += processNode(c)
				}
			case "sup":
				switch extractClass(c) {
				case "superscript":
					// 未対応
					isProcessed = true
					result += "(" + processNode(c) + ")"
				case "okurigana":
					// 未対応
					isProcessed = true
					result += processNode(c)
				}
			case "em":
				// 傍線・傍点はイタリック体で表現
				result += style(processNode(c), styleItalic)
				isProcessed = true
			case "a":
				if extractClass(c) == "midashi_anchor" {
					isProcessed = true
					result += processNode(c)
				}
			case "img":
				switch extractClass(c) {
				case "gaiji":
					isProcessed = true
					result += "\\_b[" + extractAttr(c, "src") + ",inline,opaque]"
				case "illustration":
					// 未対応
					isProcessed = true
					result += "〔画像〕"
				}
			case "h1":
				// タイトル
				isProcessed = true
				result += style(processNode(c), styleBold)
			case "h2":
				// 作者
				isProcessed = true
				result += style(processNode(c), styleBold)
			case "h3":
				// 大見出し
				isProcessed = true
				result += style(processNode(c), styleBold)
			case "h4":
				// 中見出し
				isProcessed = true
				result += style(processNode(c), styleBold)
			case "h5":
				// 小見出し
				isProcessed = true
				result += style(processNode(c), styleBold)
			case "br":
				// 改行
				isProcessed = true
				result += setting.ReturnCode
			case "ruby":
				// ルビ付き文字
				isProcessed = true
				result += expandRuby(c)
			}
			if !isProcessed {
				// その他の属性タグはここで処理
				text := c.Data
				if class := extractClass(c); class != "" {
					text += "." + class
				}
				text = "{" + text + "}"
				result += text + processNode(c)
			}
		} else if c.Type == html.TextNode {
			text := strings.ReplaceAll(c.Data, "\n", "")
			if text != "" {
				result += text
			}
		}
	}
	return replastCRLF.ReplaceAllString(result, "")
}
