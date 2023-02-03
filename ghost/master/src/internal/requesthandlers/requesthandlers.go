package requesthandlers

import (
	shiori "github.com/Narazaka/shiorigo"
	"github.com/apxxxxxxe/gohst/dictionary"
)

/* main パッケージから適宜 OnLoad(), OnUnload(), OnRequest() が呼ばれます。
 *
 * OnRequest() が呼ばれると、リクエストのイベント ID に応じて Handlers というハンドラテーブルから対応したハンドラを探して実行します。
 * ハンドラは OnLoad() が呼ばれた時に登録しておきます。
 * 各ハンドラではリクエストを基にレスポンスを作成して返します。
 * 特に返す内容が無かったり、リクエストを無視するときは 204 No Content を返します。
 * これにより「正常に終了したが返すものはない」という意図を伝えることができます。
 *
 * 他に、レスポンスの作成が簡単にできるような機能も定義しています。
 *
 * 例えば、Talks 型でトークのリストを定義してその中からランダムにひとつ選んで返すことができます。
 * このとき、各トークはバックスラッシュをエスケープした Sakura Script です。
 *
 * 各 Response〜 系の関数は特定のレスポンスを返すためのレスポンスビルダです。
 * さらに CreateGetHandlerOf() を使うと GET イベントの時に特定の値を返すハンドラを作成できます。
 *
 * ランダムトークは OnSecondChange イベントの度にカウンタを増やし、閾値を超えたら 1/10 の確率でトークを返すという実装になっています。
 * 閾値を超えたがトークを返さない場合に、乱数のシードを変えています。
 * これにより繰り返しがよりランダムに起こることを期待しています。
 *
 * なお、このパッケージでは擬似乱数発生に rand パッケージを使っています。
 * この乱数のシードの更新は、上述の OnSecondChange イベントと OnLoad() が呼ばれた時に ResetRNG() を呼び出して行っています。
 *
 * このパッケージではトークを文字リテラルでコード中に埋め込んでいますが、TOML や Lua などの形で外部にデータを移すとゴースト開発が楽になるでしょう。
 * また、NOTIFY リクエストなどで得た情報を変数に記録しておき、text/template でトークに埋め込むことも可能です。
 * 単語をその種類ごとに分けた辞書を map[string][]string で作り、ある種類の単語を辞書からランダムに取ってくるのは、トークの多様性を増すための古典的な方法です。
 *
 * イベントの種類など SHIORI/3.0 の仕様については「http://ssp.shillest.net/ukadoc/manual/index.html」を見てください。
 */

var (
	Dictionary *dictionary.Dictionary
)

func OnLoad(path string) error {
	var err error
	// 辞書と変数をDictionaryに読み込む
	Dictionary, err = dictionary.New()
	if err != nil {
		return err
	}

	return nil
}

func OnUnload() error {
	// 変数をセーブする
	if err := dictionary.SaveVariables(*Dictionary.Variables); err != nil {
		return err
	}

	return nil
}

func OnRequest(req shiori.Request) (shiori.Response, error) {
	// リクエストヘッダがなければ、統一的操作のために初期化しておきます。
	if req.Headers == nil {
		req.Headers = shiori.RequestHeaders{}
	}

	// ID ヘッダにはイベント名が入っています。
	// イベントに対応するハンドラが定義されていれば呼び出します。
	if event, ok := req.Headers["ID"]; ok {
		if handler, ok := Dictionary.Handlers[event]; ok {
			return handler(req, Dictionary.Variables)
		}
	}

	return dictionary.ResponseNoContent(), nil
}
