package save

import (
	"fmt"
	"os"

	"golang.org/x/xerrors"

	"ikimonoaz-exporter/template"
	"ikimonoaz-exporter/userdata"
)

func prepareDirectories(dir string) error {
	// ディレクトリを用意
	if err := os.Mkdir(dir+"articles", 0755); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	if err := os.Mkdir(dir+"media", 0755); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	return nil
}

func saveIndex(dir string, ud userdata.UserData) error {
	context := ud.ToMap()
	s, err := template.RenderIndex(context)
	if err != nil {
		// テンプレートへのレンダリングに失敗した
		return err
	}

	path := fmt.Sprintf("%sindex.html", dir)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	f.Truncate(0)
	f.Seek(0, 0)
	if _, err := f.WriteString(s); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func saveArticle(dir string, article userdata.Article) error {
	context := article.ToMap()
	s, err := template.RenderArticle(context)
	if err != nil {
		// テンプレートへのレンダリングに失敗した
		return err
	}

	path := fmt.Sprintf("%sarticles/%d.html", dir, article.ID)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	f.Truncate(0)
	f.Seek(0, 0)
	if _, err := f.WriteString(s); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func SaveUserData(dir string, ud userdata.UserData) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return xerrors.Errorf("保存先フォルダが存在しません")
	}

	fmt.Println("[ikimonoaz-exporter] 保存ディレクトリ用意中...")
	if err := prepareDirectories(dir); err != nil {
		fmt.Printf("%v\n", err)
		return xerrors.Errorf("データ保存用フォルダの作成に失敗しました")
	}

	fmt.Println("[ikimonoaz-exporter] indexページ保存中...")
	if err := saveIndex(dir, ud); err != nil {
		fmt.Printf("%v\n", err)
		return xerrors.Errorf("ユーザ情報のエクスポートに失敗しました")
	}

	fmt.Println("[ikimonoaz-exporter] 記事ページ保存中...")
	for _, a := range ud.Articles {
		fmt.Printf("[ikimonoaz-exporter] 記事ID: %d\n", a.ID)

		// デバッグ用: 最初の1記事のメディアデータだけ取得する
		// if i >= 1 {
		// 	break
		// }

		fmt.Println("[ikimonoaz-exporter] 記事ページ保存中...")
		if err := saveArticle(dir, a); err != nil {
			fmt.Printf("%v\n", err)
			return xerrors.Errorf("記事のエクスポートに失敗しました")
		}

		fmt.Println("[ikimonoaz-exporter] 記事のメディアファイル保存中...")
		for _, m := range a.MediaList {
			if err := saveMedia(dir, m.URL); err != nil {
				fmt.Printf("%v\n", err)
				return xerrors.Errorf("画像・動画のダウンロードに失敗しました")
			}
		}
	}

	return nil
}
