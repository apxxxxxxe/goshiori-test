package balloon

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/apxxxxxxe/goshiori-test/aozoragetter/setting"
)

const punctuations = "、。)）」』}｝]〕】》>＞〉"

var tagRep = regexp.MustCompile(`\\_{0,2}[a-zA-Z0-9*!&\-+](\d|\[("([^"]|\\")+?"|([^\]]|\\\])+?)+?\])?`)
var verticalUnicode = map[string]string{
	",": string(rune(0xFE10)),
	"，": string(rune(0xFE10)),
	"、": string(rune(0xFE11)),
	".": string(rune(0xFE12)),
	"．": string(rune(0xFE12)),
	"。": string(rune(0xFE12)),
	":": string(rune(0xFE13)),
	"：": string(rune(0xFE13)),
	";": string(rune(0xFE14)),
	"；": string(rune(0xFE14)),
	"!": string(rune(0xFE15)),
	"！": string(rune(0xFE15)),
	"?": string(rune(0xFE16)),
	"？": string(rune(0xFE16)),
	"…": string(rune(0xFE19)),
	"ー": "｜",
	"−": string(rune(0x2758)),
	"(": string(rune(0xFE35)),
	"（": string(rune(0xFE35)),
	")": string(rune(0xFE36)),
	"）": string(rune(0xFE36)),
	"「": string(rune(0xFE41)),
	"」": string(rune(0xFE42)),
	"『": string(rune(0xFE43)),
	"』": string(rune(0xFE44)),
	"{": string(rune(0xFE37)),
	"｛": string(rune(0xFE37)),
	"}": string(rune(0xFE38)),
	"｝": string(rune(0xFE38)),
	"[": string(rune(0xFE39)),
	"〔": string(rune(0xFE39)),
	"]": string(rune(0xFE3A)),
	"〕": string(rune(0xFE3A)),
	"【": string(rune(0xFE3B)),
	"】": string(rune(0xFE3C)),
	"《": string(rune(0xFE3D)),
	"》": string(rune(0xFE3E)),
	"<": string(rune(0xFE3F)),
	"＜": string(rune(0xFE3F)),
	"〈": string(rune(0xFE3F)),
	">": string(rune(0xFE40)),
	"＞": string(rune(0xFE40)),
	"〉": string(rune(0xFE40)),
	"〜": string(rune(0x2307)),
}

type Balloon struct {
	Width     int
	Height    int
	FontSize  int
	FontName  string
	Margin    int
	BaseYPos  int
	XPos      int
	YPos      int
	PageIndex int
	Pages     []string
}

func (b *Balloon) InitialScripts() string {
	res := "\\b[2]\\f[name," + b.FontName + "]\\f[height," + strconv.Itoa(b.FontSize) + "]"

	pageNum := strconv.Itoa(b.PageIndex + 1)
	xPos := (b.Width - (b.FontSize / 2 * len(pageNum))) / 2
	yPos := (b.Height - b.BaseYPos) - b.FontSize + b.Height*b.PageIndex - b.FontSize/2
	res += "\\_l[" + strconv.Itoa(xPos) + "," + strconv.Itoa(yPos) + "]\\_q" + pageNum + "\\_q"

	res += "\\![set,autoscroll,disable]"

	return res
}

// バルーンの書き込み位置情報をリセットする
func (b *Balloon) Reset() {
	b.XPos = 0
	b.YPos = 0
	b.PageIndex = 0
	b.Pages = []string{b.InitialScripts()}
}

// srcをバルーン内バッファに縦書き変換しつつ記憶し、その内容分だけ書き込み位置を進める
func (b *Balloon) StoreVirticalWriting(src string) {
	tagMark := string(rune(0x2))
	tags := tagRep.FindAllString(src, -1)
	src = tagRep.ReplaceAllString(src, tagMark)

	tagCount := 0
	words := strings.Split(src, "")
	for i, w := range words {
		if w == tagMark {
			// さくらスクリプト
			tag := tags[tagCount]
			switch tag {
			case setting.ReturnCode:
				b.Return()
			case setting.RepageCode:
				b.RePage()
			default:
				b.Pages[b.PageIndex] += tag
			}
			tagCount++
		} else {
			// 通常の文字
			if v, ok := verticalUnicode[w]; ok {
				// 縦書き文字があるなら変換
				if v == string(rune(0xFE35)) || v == string(rune(0xFE36)) {
					w = v + "\\_l[@-500em,]" + w
				} else {
					w = "\\_q" + v + "\\_q\\_l[@-500em,]" + w
				}
			}
			wordMax := (b.Height-b.BaseYPos)/b.FontSize - 2
			xPos := (b.Width - b.FontSize) - (b.FontSize+b.Margin)*b.XPos
			yPos := b.YPos*b.FontSize + b.Height*b.PageIndex
			b.Pages[b.PageIndex] += "\\_l[" + strconv.Itoa(xPos) + "," + strconv.Itoa(yPos) + "]" + w
			b.YPos++
			if b.YPos >= wordMax && (i == len(words)-1 || !strings.ContainsAny(words[i+1], punctuations)) {
				// 改行処理
				b.Return()
			}
		}
	}
}

// 改行処理
func (b *Balloon) Return() {
	lineMax := b.Width / (b.Margin + b.FontSize)
	b.YPos = 0
	b.XPos++
	if b.XPos == lineMax {
		b.RePage()
	}
}

// 改ページ処理
func (b *Balloon) RePage() {
	b.PageIndex++
	b.XPos = 0
	b.YPos = 0
	b.Pages = append(b.Pages, b.InitialScripts())
}

// Pagesの内容をファイルに出力
func (b *Balloon) WriteFile(destPath string) error {
	fp, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer fp.Close()

	for _, p := range b.Pages {
		if _, err := fp.WriteString(p + "\n"); err != nil {
			return err
		}
	}
	return nil
}
