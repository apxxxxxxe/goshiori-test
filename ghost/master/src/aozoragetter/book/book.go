package book

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/apxxxxxxe/goshiori-test/aozoragetter/format"
	"github.com/apxxxxxxe/goshiori-test/aozoragetter/setting"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func filterBook(queries []string, indexData []Book) []Book {
	candidates := indexData
	for _, q := range queries {
		candidates = searchBook(q, candidates)
	}
	return candidates
}

func searchBook(query string, candidates []Book) []Book {
	results := []Book{}
	for _, book := range candidates {
		src := fmt.Sprintf("%s %s %s %s", book.Title, book.Author, book.ID, book.Mojidukai)
		if strings.Contains(src, query) {
			results = append(results, book)
		}
	}
	return results
}

func fetchBookFile(book Book, path string) (string, error) {
	tmpPath := filepath.Join(path, "tmp")
	if err := setting.DownloadFile(tmpPath, book.HtmlURL); err != nil {
		return "", err
	}

	fp, err := os.Open(tmpPath)
	if err != nil {
		return "", err
	}
	defer func() {
		fp.Close()
		os.RemoveAll(tmpPath)
	}()

	var r io.Reader
	if book.HtmlEncoding == "ShiftJIS" {
		r = transform.NewReader(fp, japanese.ShiftJIS.NewDecoder())
	} else {
		r = fp
	}
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func GetBook(queries []string) ([]Book, string, error) {
	if len(queries) == 0 {
		return nil, "", nil
	}

	exec, err := os.Executable()
	if err != nil {
		return nil, ",", err
	}
	dir := filepath.Dir(exec)

	index, err := MakeIndex(dir)
	if err != nil {
		return nil, "", err
	}

	candidates := filterBook(queries, index)

	switch len(candidates) {
	case 0:
		return nil, "", err
	case 1:
		book := candidates[0]
		rawText, err := fetchBookFile(book, dir)
		if err != nil {
			return nil, "", err
		}
		text, err := format.ParseHtmlText(rawText)
		return candidates, text, err
	default:
		return candidates, "", nil
	}
}
