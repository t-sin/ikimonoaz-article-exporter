package main

import (
	// "encoding/json"
	// "fmt"
	// "io/ioutil"
	// "os"

	// "ikimonoaz-exporter/ikimonoaz"
	// "ikimonoaz-exporter/save"
	// "ikimonoaz-exporter/userdata"

	"ikimonoaz-exporter/gui"
)

func export(targetPath, mypageURL string) error {
	return nil
}

func main() {
	// userID := "7308"
	// userdata, err := ikimonoaz.CollectAllUserData(userID)
	// if err != nil {
	// 	// データ (メディアファイル以外) 収集に失敗した
	// 	fmt.Printf("%+v\n", err)
	// 	return
	// }

	// // デバッグ用: ユーザデータを取得したら標準出力に書き出す
	// u, _ := json.Marshal(&userdata)
	// fmt.Printf("%s\n", u)

	// デバッグ用: 標準入力からユーザデータを入れるとそれを元にエクスポートする
	// bytes, _ := ioutil.ReadAll(os.Stdin)
	// var userdata userdata.UserData
	// json.Unmarshal(bytes, &userdata)

	// targetDir := "./testdir/"

	// if err := save.SaveUserData(targetDir, userdata); err != nil {
	// 	fmt.Printf("%+v\n", err)
	// }

	gui.Start(export)
}
