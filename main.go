package main

import (
	"fmt"

	"ikimonoaz-exporter/ikimonoaz"
)

func main() {
	userID := "7308"
	articles, err := ikimonoaz.CollectAllUserData(userID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return
	}
	fmt.Printf("%v\n", articles)
}
