package save

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/xerrors"

	"ikimonoaz-exporter/template"
	"ikimonoaz-exporter/userdata"
)

func prepareDirectories(dir string) error {
	// ディレクトリを用意
	path := fmt.Sprintf("%s%c%s", dir, os.PathSeparator, "articles")
	if err := os.Mkdir(path, 0755); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	path = fmt.Sprintf("%s%c%s", dir, os.PathSeparator, "media")
	if err := os.Mkdir(path, 0755); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	return nil
}

func saveJSON(dir string, ud userdata.UserData) error {
	u, err := json.Marshal(ud)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s%cuserdata.json", dir, os.PathSeparator)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	f.Truncate(0)
	f.Seek(0, 0)
	if _, err := f.WriteString(string(u)); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
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

	path := fmt.Sprintf("%s%cindex.html", dir, os.PathSeparator)
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

	path := fmt.Sprintf("%s%carticles%c%d.html", dir, os.PathSeparator, os.PathSeparator, article.ID)
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

	fmt.Println("[ikimonoaz-exporter] JSONデータ保存中...")
	if err := saveJSON(dir, ud); err != nil {
		fmt.Printf("%v\n", err)
		return xerrors.Errorf("ユーザ情報JSONのエクスポートに失敗しました")
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
