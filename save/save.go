package save

import (
	"fmt"
	"os"

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
	if err := prepareDirectories(dir); err != nil {
		return err
	}

	if err := saveIndex(dir, ud); err != nil {
		return err
	}

	for _, a := range ud.Articles {
		if err := saveArticle(dir, a); err != nil {
			return err
		}

		for _, m := range a.MediaList {
			if err := saveMedia(dir, m.URL); err != nil {
				return err
			}
		}
	}

	return nil
}
