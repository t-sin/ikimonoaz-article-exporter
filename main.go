package main

import (
	"fmt"

	"ikimonoaz-exporter/ikimonoaz"
	"ikimonoaz-exporter/template"
)

func main() {
	userID := "7308"
	userdata, err := ikimonoaz.CollectAllUserData(userID)
	if err != nil {
		// データ (メディアファイル以外) 収集に失敗した
		fmt.Printf("%+v\n", err)
		return
	}

	context := userdata.ToMap()
	s, err := template.RenderIndex(context)
	if err != nil {
		// テンプレートへのレンダリングに失敗した
		fmt.Printf("%+v\n", err)
		return
	}

	fmt.Println(s)
}
