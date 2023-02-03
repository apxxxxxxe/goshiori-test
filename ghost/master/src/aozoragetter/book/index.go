package book

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/apxxxxxxe/goshiori-test/aozoragetter/setting"
)

type Book struct {
	ID           string
	Title        string
	Author       string
	HtmlURL      string
	HtmlEncoding string
	Mojidukai    string
}

const indexFile = "list_person_all_extended_utf8.csv"

// ローカルにindexFileを取得する
func getIndexCSV(baseDir string) error {
	indexZip := filepath.Join(baseDir, "tmp.zip")
	if !setting.IsFile(filepath.Join(baseDir, indexFile)) {
		if err := setting.DownloadFile(indexZip, "https://www.aozora.gr.jp/index_pages/list_person_all_extended_utf8.zip"); err != nil {
			return err
		}

		if err := setting.Unzip(indexZip, baseDir); err != nil {
			return err
		}

		if err := os.RemoveAll(indexZip); err != nil {
			return err
		}
	}

	return nil
}

// indexFileを配列に読み込んで返す
func loadIndexCSV(path string) ([][]string, error) {
	const delim = ','

	s, err := os.ReadFile(path)
	if err != nil {
		return [][]string{}, err
	}

	r := csv.NewReader(strings.NewReader(string(s)))
	r.Comma = delim

	result, err := r.ReadAll()
	if err != nil {
		return [][]string{}, err
	}

	return result, nil
}

func isKatakana(src string) bool {
	result := true
	for _, r := range src {
		if !unicode.In(rune(r), unicode.Katakana) {
			result = false
		}
	}
	return result
}

// indexFile一行分の情報をBook構造体に格納して返す
func getInfoSummury(bookInfo []string) Book {
	result := Book{}
	result.ID = bookInfo[0]
	result.Title = bookInfo[1]
	result.HtmlURL = bookInfo[50]
	result.HtmlEncoding = bookInfo[52]
	result.Mojidukai = bookInfo[9]
	if (isKatakana(bookInfo[15]) && isKatakana(bookInfo[16])) || strings.Contains(bookInfo[16], "・") {
		result.Author = bookInfo[16] + "・" + bookInfo[15]
	} else {
		result.Author = bookInfo[15] + bookInfo[16]
	}
	return result
}

// indexFileの情報をすべてBook構造体に格納して返す
func MakeIndex(dirname string) ([]Book, error) {
	if err := os.MkdirAll(dirname, 0755); err != nil {
		return nil, err
	}

	if err := getIndexCSV(dirname); err != nil {
		return nil, err
	}

	csv, err := loadIndexCSV(filepath.Join(dirname, indexFile))
	if err != nil {
		return nil, err
	}

	result := []Book{}
	for _, info := range csv {
		result = append(result, getInfoSummury(info))
	}

	return result, nil
}
