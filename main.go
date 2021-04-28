package main

import (
	// "encoding/json"
	"fmt"
	// "io/ioutil"
	"os"
	"regexp"

	"golang.org/x/xerrors"

	"ikimonoaz-exporter/gui"
	"ikimonoaz-exporter/ikimonoaz"
	"ikimonoaz-exporter/save"
	// "ikimonoaz-exporter/userdata"
)

var mypageURLPat = regexp.MustCompile(`^https://ikimonoaz.ikimonopal.jp/profile/u/([0-9]+)$`)

func export(targetPath, mypageURL string) error {
	fmt.Println("[ikimonoaz-exporter] エクスポート開始")

	if targetPath == "" {
		return xerrors.Errorf("保存先フォルダを選んでください")
	}

	if os.IsPathSeparator(targetPath[len(targetPath)-1]) {
		targetPath = targetPath + os.PathSeparator
	}

	matches := mypageURLPat.FindStringSubmatch(mypageURL)
	if len(matches) != 2 {
		return xerrors.Errorf("マイページのURLでない文字列が入力されました")
	}

	fmt.Println("[ikimonoaz-exporter] ユーザ情報および記事データ取得中...")
	userID := matches[1]
	userdata, err := ikimonoaz.CollectAllUserData(userID)
	if err != nil {
		return xerrors.Errorf("記事データの収集に失敗しました。\n時間をおいてもダメな場合作者に連絡ください。")
	}

	fmt.Println("[ikimonoaz-exporter] メディアデータ取得中...")
	if err := save.SaveUserData(targetPath, userdata); err != nil {
		return err
	}

	fmt.Println("[ikimonoaz-exporter] エクスポート完了")

	return nil
}

func main() {
	// デバッグ用: 標準入力からユーザデータを入れるとそれを元にエクスポートする
	// bytes, _ := ioutil.ReadAll(os.Stdin)
	// var userdata userdata.UserData
	// json.Unmarshal(bytes, &userdata)

	gui.Start(export)
}
